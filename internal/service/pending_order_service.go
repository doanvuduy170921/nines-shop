package service

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"nineshop-be/internal/worker"
	"nineshop-be/pkg/cache"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"nineshop-be/internal/constant"
	"nineshop-be/internal/db"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/email"
	"nineshop-be/pkg/vnpay"
)

type pendingOrderService struct {
	store       *db.Store
	repo        repository.PendingOrderRepository
	ui          UserInject
	pot         repository.PendingOrderItemRepository
	email       email.EmailService
	or          repository.OrderRepository
	oi          repository.OrderItemRepository
	pr          repository.ProductRepository
	osh         repository.OrderStatusHistoryRepository
	emailWorker worker.BackgroundWorker
	redis       cache.RedisCacheService
}

type UserInject interface {
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
}

func NewPendingOrderService(emailWorker worker.BackgroundWorker, store *db.Store, repo repository.PendingOrderRepository, ui UserInject, pot repository.PendingOrderItemRepository, email email.EmailService, or repository.OrderRepository, oi repository.OrderItemRepository, pr repository.ProductRepository, osh repository.OrderStatusHistoryRepository, redis cache.RedisCacheService) PendingOrderService {
	return &pendingOrderService{
		store:       store,
		repo:        repo,
		ui:          ui,
		pot:         pot,
		email:       email,
		or:          or,
		oi:          oi,
		pr:          pr,
		osh:         osh,
		emailWorker: emailWorker,
		redis:       redis,
	}
}

type orderTxResult struct {
	PendingOrder   sqlc.PendingOrder
	OTP            string
	ShouldSendOTP  bool
	VnPayRequested bool
}

func (ns *pendingOrderService) Create(ctx *gin.Context, arg dto.CreatePendingOrderDto) (dto.CreatePendingOrderResponse, error) {
	c := ctx.Request.Context()
	user, err := ns.extractAndValidateUser(ctx)
	if err != nil {
		log.Printf("extract and validater user error : %v", err)
		return dto.CreatePendingOrderResponse{}, err
	}

	if len(arg.Items) == 0 {
		log.Printf("arg.Items is empty")
		return dto.CreatePendingOrderResponse{}, utils.NewError(400, "items cannot be empty")
	}
	txResult, err := ns.executeOrderTransaction(c, user, arg)
	if err != nil {
		log.Printf("executeOrderTransaction error : %v", err)
		return dto.CreatePendingOrderResponse{}, err
	}

	paymentURL, err := ns.handlePostTransactionActions(ctx, txResult, user.Email)
	if err != nil {
		log.Printf("handlePostTransactionActions error : %v", err)
		return dto.CreatePendingOrderResponse{}, err
	}
	po := txResult.PendingOrder

	return dto.CreatePendingOrderResponse{
		ID:              po.ID,
		UserID:          user.ID,
		Name:            po.Name,
		Email:           po.Email,
		Phone:           po.Phone,
		PaymentMethodID: po.PaymentMethodID,
		ShippingPrice:   po.ShippingPrice,
		Subtotal:        po.Subtotal,
		TotalAmount:     po.TotalAmount,
		CreatedAt:       po.CreatedAt,
		Status:          po.Status,
		Tax:             po.Tax,
		AmountItem:      po.AmountItem,
		Address:         po.Address,
		VnPayURL:        paymentURL,
		ExpiresAt:       po.ExpiresAt,
	}, nil
}

func (ns *pendingOrderService) extractAndValidateUser(ctx *gin.Context) (sqlc.User, error) {
	userUuid, exists := ctx.Get("user_uuid")
	if !exists {
		log.Printf("userUuid is empty")
		return sqlc.User{}, utils.NewError(400, "user uuid not found in request context")
	}

	uuidStr, ok := userUuid.(string)
	if !ok {
		log.Printf("uuidStr is empty")
		return sqlc.User{}, utils.NewError(400, "user uuid format is invalid")
	}

	pgUuid, err := utils.StringToPgUuid(uuidStr)
	if err != nil {
		log.Printf("uuidStr is invalid")
		return sqlc.User{}, utils.NewError(400, "convert user uuid to uuid fail")
	}

	c := ctx.Request.Context()
	user, err := ns.store.GetByUuid(c, pgUuid)
	if err != nil {
		log.Printf("get user by uuid fail: %v", err)
		return sqlc.User{}, utils.NewError(http.StatusBadRequest, "get user by uuid fail")
	}

	return user, nil
}

func (ns *pendingOrderService) executeOrderTransaction(c context.Context, user sqlc.User, arg dto.CreatePendingOrderDto) (orderTxResult, error) {
	var pOrder sqlc.PendingOrder
	var otp string
	var otpExpiresAt pgtype.Timestamp
	var shouldSendOTP bool
	var vnPayRequested bool
	now := time.Now().UTC()
	var expiresAt time.Time

	err := ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		items := slices.Clone(arg.Items)
		sort.Slice(items, func(i, j int) bool { return items[i].VariantID < items[j].VariantID }) // sort variant for fix deadlock
		for _, item := range items {
			variant, errTx := q.GetVariantForUpdate(c, int32(item.VariantID))
			if errTx != nil {
				log.Printf("get variant fail: %v", errTx)
				return utils.WrapError(errTx, "Product variant not found", http.StatusBadRequest)
			}

			if !(*variant.IsActive) {
				log.Printf("variant is not active")
				return utils.NewError(http.StatusBadRequest, "Product is inactive")
			}
			if variant.StockQuantity < int32(item.Quantity) {
				log.Printf("variant stock quantity is less then stock quantity")
				return utils.NewError(http.StatusBadRequest, "Insufficient stock quantity")
			}

			result, errTx := q.DecreaseStock(c, sqlc.DecreaseStockParams{
				Stock:     int32(item.Quantity),
				VariantID: int32(item.VariantID),
			})
			if errTx != nil || result.RowsAffected() == 0 {
				log.Printf("product stock fail: %v", errTx)
				return utils.NewError(http.StatusBadRequest, "Product is out of stock or unavailable")
			}
		}

		// Handle Payment Method / OTP
		if arg.PaymentId == constant.PaymentMethodCod {
			otp = utils.GenOTP()
			otpExpiresAt = pgtype.Timestamp{Time: now.Add(5 * time.Minute), Valid: true}
			expiresAt = now.Add(5 * time.Minute)
			shouldSendOTP = true
		} else {
			otp = ""
			otpExpiresAt = pgtype.Timestamp{Valid: false}
			expiresAt = now.Add(15 * time.Minute)
			vnPayRequested = arg.PaymentId == constant.PaymentMethodVnPay
		}

		count := int32(len(arg.Items))
		var errTx error

		pOrder, errTx = q.CreatePendingOrder(c, sqlc.CreatePendingOrderParams{
			UserID:          user.ID,
			Name:            arg.Name,
			Email:           user.Email,
			Phone:           arg.Phone,
			PaymentMethodID: int32(arg.PaymentId),
			Address:         arg.Address,
			Otp:             otp,
			OtpExpiresAt:    otpExpiresAt,
			ExpiresAt:       pgtype.Timestamptz{Time: expiresAt, Valid: true},
			Subtotal:        utils.Float64ToNumeric(arg.Subtotal),
			TotalAmount:     utils.Float64ToNumeric(arg.Total),
			Tax:             utils.Float64ToNumeric(arg.Tax),
			ShippingPrice:   utils.Float64ToNumeric(arg.Shipping),
			AmountItem:      &count,
		})
		if errTx != nil {
			log.Printf("ERROR creating pending order: %v", errTx)
			return utils.WrapError(errTx, "Create Pending order fail", http.StatusBadRequest)
		}
		for _, item := range arg.Items {
			_, errTx = q.CreatePendingOrderItem(c, sqlc.CreatePendingOrderItemParams{
				PendingOrderID: pOrder.ID,
				VariantID:      int32(item.VariantID),
				Quantity:       int32(item.Quantity),
				Price:          utils.Float64ToNumeric(item.Price),
			})
			if errTx != nil {
				log.Printf("ERROR creating pending order item: %v", errTx)
				return utils.WrapError(errTx, "Create item fail", http.StatusBadRequest)
			}
		}

		return nil
	})

	if err != nil {
		return orderTxResult{}, err
	}

	return orderTxResult{
		PendingOrder:   pOrder,
		OTP:            otp,
		ShouldSendOTP:  shouldSendOTP,
		VnPayRequested: vnPayRequested,
	}, nil
}

func (ns *pendingOrderService) handlePostTransactionActions(ctx *gin.Context, txResult orderTxResult, userEmail string) (string, error) {
	var paymentURL string
	var err error

	if txResult.VnPayRequested {
		totalAmountNumeric, _ := txResult.PendingOrder.TotalAmount.Float64Value()
		vnPayService := vnpay.NewService()

		paymentURL, err = vnPayService.CreatePaymentURL(vnpay.OrderInfo{
			ID:          int64(txResult.PendingOrder.ID),
			TotalAmount: totalAmountNumeric.Float64,
		}, ctx.ClientIP())
		if err != nil {
			log.Printf("ERROR creating VNPAY URL: %v", err)
			return "", utils.WrapError(err, "Create VNPAY URL fail", http.StatusInternalServerError)
		}
	}

	if txResult.ShouldSendOTP {
		ns.emailWorker.EnqueueEmailTask(userEmail, txResult.OTP)
		log.Printf("Moved the OTP email-sending task to a background worker for user: %s", userEmail)
	}

	return paymentURL, nil
}

const maxOtpAttempts = 5

func (ns *pendingOrderService) releaseAndLog(c context.Context, id int32, status string) {
	if err := ns.releasePending(c, id, status); err != nil {
		log.Printf("release pending order %d (%s) fail: %v", id, status, err)
	}
}

func (ns *pendingOrderService) ValidateOTP(ctx *gin.Context, arg dto.ValidateOTPParams) (sqlc.Order, error) {
	c := ctx.Request.Context()

	user, err := ns.extractAndValidateUser(ctx)
	if err != nil {
		return sqlc.Order{}, err
	}

	pOrder, err := ns.repo.GetPOrderById(c, int32(arg.POrderId))
	if err != nil || pOrder.UserID != user.ID {
		return sqlc.Order{}, utils.NewError(http.StatusNotFound, "Order not found")
	}
	if pOrder.Status != constant.PendingOrderStatusPending {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "The order is no longer valid")
	}

	if time.Now().After(pOrder.ExpiresAt.Time) {
		ns.releaseAndLog(c, pOrder.ID, constant.PendingOrderStatusExpired)
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "The order has expired")
	}

	if pOrder.OtpAttempts >= maxOtpAttempts {
		ns.releaseAndLog(c, pOrder.ID, constant.PendingOrderStatusCancelled)
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "You have entered an incorrect OTP more than 5 times. The order has been cancelled.")
	}
	if !pOrder.OtpExpiresAt.Valid || time.Now().After(pOrder.OtpExpiresAt.Time) {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "OTP has expired")
	}

	if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(arg.Otp)), []byte(pOrder.Otp)) != 1 {
		attempts, err := ns.store.IncrementOtpAttempts(c, pOrder.ID)
		if err != nil {
			return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "The order is no longer valid")
		}
		if attempts >= maxOtpAttempts {
			ns.releaseAndLog(c, pOrder.ID, constant.PendingOrderStatusCancelled)
			return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "You have entered an incorrect OTP more than 5 times. The order has been cancelled")
		}
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest,
			fmt.Sprintf("The OTP is incorrect. You have %d attempt(s) remaining", maxOtpAttempts-attempts))
	}

	var order sqlc.Order
	err = ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		if _, err := q.ConfirmPendingOrder(c, sqlc.ConfirmPendingOrderParams{ID: pOrder.ID, UserID: user.ID}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return utils.NewError(http.StatusBadRequest, "The order is no longer valid")
			}
			return utils.WrapError(err, "Confirm pending order fail", http.StatusInternalServerError)
		}

		var errTx error
		order, errTx = q.CreateOrder(c, sqlc.CreateOrderParams{
			UserID:          pOrder.UserID,
			Name:            pOrder.Name,
			Email:           pOrder.Email,
			Phone:           pOrder.Phone,
			Address:         pOrder.Address,
			PaymentMethodID: pOrder.PaymentMethodID,
			ShippingPrice:   pOrder.ShippingPrice,
			Subtotal:        pOrder.Subtotal,
			TotalAmount:     pOrder.TotalAmount,
			Tax:             pOrder.Tax,
			Status:          constant.OrderStatusPending,
			AmountItem:      pOrder.AmountItem,
			PaymentStatus:   constant.PaymentStatusUnpaid,
		})
		if errTx != nil {
			return utils.WrapError(errTx, "Create order fail", http.StatusInternalServerError)
		}

		statusHistory := constant.OrderStatusPending
		note := constant.OrderStatusNotes[statusHistory]
		_, errTx = q.CreateOrderStatusHistory(c, sqlc.CreateOrderStatusHistoryParams{
			OrderID: order.ID,
			Status:  statusHistory,
			Note:    &note,
		})
		if errTx != nil {
			return utils.WrapError(errTx, "Create order status history fail", http.StatusInternalServerError)
		}

		pOrderItems, errTx := q.GetByPOrderItemId(c, pOrder.ID)
		if errTx != nil {
			return utils.WrapError(errTx, "Get pending order items fail", http.StatusInternalServerError)
		}
		for _, item := range pOrderItems {
			product, errVariant := q.GetVariantById(c, item.VariantID)
			if errVariant != nil {
				return utils.WrapError(errVariant, "Get product from pending order item fail", http.StatusBadRequest)
			}
			_, errTx = q.AddOrderItem(c, sqlc.AddOrderItemParams{
				OrderID:          order.ID,
				VariantID:        item.VariantID,
				Quantity:         item.Quantity,
				Price:            item.Price,
				ProductName:      product.ProductName,
				ProductThumbnail: &product.ProductThumbnail,
			})
			if errTx != nil {
				return utils.WrapError(errTx, "Add order item fail", http.StatusInternalServerError)
			}
		}
		return nil
	})
	if err != nil {
		return sqlc.Order{}, err
	}
	return order, nil
}

func (ns *pendingOrderService) releasePending(c context.Context, id int32, newStatus string) error {
	return ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		_, err := q.UpdatePendingOrderStatus(c, sqlc.UpdatePendingOrderStatusParams{Status: newStatus, ID: id})
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		items, err := q.GetByPOrderItemId(c, id)
		if err != nil {
			return err
		}
		sort.Slice(items, func(i, j int) bool { return items[i].VariantID < items[j].VariantID })
		for _, it := range items {
			if err := q.IncreaseStock(c, sqlc.IncreaseStockParams{Stock: it.Quantity, VariantID: it.VariantID}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ns *pendingOrderService) ReleaseExpiredPendingOrders(ctx context.Context) (int, error) {
	ids, err := ns.store.ListExpiredPendingOrderIDs(ctx)
	if err != nil {
		return 0, err
	}
	released := 0
	for _, id := range ids {
		if err := ns.releasePending(ctx, id, constant.PendingOrderStatusExpired); err != nil {
			log.Printf("release expired pending order %d fail: %v", id, err)
			continue
		}
		released++
	}
	return released, nil
}

func StartExpiredPendingOrderJob(ctx context.Context, svc PendingOrderService, interval time.Duration) {
	run := func() {
		for {
			n, err := svc.ReleaseExpiredPendingOrders(ctx)
			if err != nil {
				log.Printf("expired pending order job error: %v", err)
				return
			}
			if n < 100 {
				return
			}
		}
	}
	go func() {
		run()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}
