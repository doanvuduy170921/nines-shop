package service

import (
	"context"
	"log"
	"net/http"
	"nineshop-be/internal/worker"
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
}

type UserInject interface {
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
}

func NewPendingOrderService(emailWorker worker.BackgroundWorker, store *db.Store, repo repository.PendingOrderRepository, ui UserInject, pot repository.PendingOrderItemRepository, email email.EmailService, or repository.OrderRepository, oi repository.OrderItemRepository, pr repository.ProductRepository, osh repository.OrderStatusHistoryRepository) PendingOrderService {
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

	paymentURL, err := ns.handlePostTransactionActions(ctx, txResult, user.Email, c)
	if err != nil {
		log.Printf("handlePostTransactionActions error : %v", err)
		return dto.CreatePendingOrderResponse{}, err
	}

	return dto.CreatePendingOrderResponse{
		PendingOrder: txResult.PendingOrder,
		VnPayURL:     paymentURL,
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
	var vnpayRequested bool
	var status = constant.OrderStatusPending

	err := ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		for _, item := range arg.Items {
			// Khóa dòng sản phẩm để kiểm tra thông tin và tồn kho
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

			// Thực hiện trừ kho
			result, errTx := q.DecreaseStock(c, sqlc.DecreaseStockParams{
				Stock:     int32(item.Quantity),
				VariantID: int32(item.VariantID),
			})
			if errTx != nil || result.RowsAffected() == 0 {
				log.Printf("product stock fail: %v", errTx)
				return utils.NewError(http.StatusBadRequest, "Product is out of stock or unavailable")
			}
		}

		// Xử lý Payment Method / OTP
		if arg.PaymentId == constant.PaymentMethodCod {
			otp = utils.GenOTP()
			otpExpiresAt = pgtype.Timestamp{
				Time:  time.Now().Add(5 * time.Minute),
				Valid: true,
			}
			shouldSendOTP = true
		} else {
			otp = ""
			otpExpiresAt = pgtype.Timestamp{Valid: false}
			vnpayRequested = arg.PaymentId == constant.PaymentMethodVnPay
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
			Subtotal:        utils.Float64ToNumeric(arg.Subtotal),
			TotalAmount:     utils.Float64ToNumeric(arg.Total),
			Tax:             utils.Float64ToNumeric(arg.Tax),
			ShippingPrice:   utils.Float64ToNumeric(arg.Shipping),
			Status:          &status,
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
		VnPayRequested: vnpayRequested,
	}, nil
}

func (ns *pendingOrderService) handlePostTransactionActions(ctx *gin.Context, txResult orderTxResult, userEmail string, c context.Context) (string, error) {
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
		log.Printf("Đã đưa task gửi email OTP vào background worker cho user: %s", userEmail)
	}

	return paymentURL, nil
}

func (ns *pendingOrderService) ValidateOTP(ctx *gin.Context, arg dto.ValidateOTPParams) (sqlc.Order, error) {
	c := ctx.Request.Context()

	var order sqlc.Order

	err := ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		pOrder, errTx := q.GetPOrderById(c, int32(arg.POrderId))
		if errTx != nil {
			return utils.WrapError(errTx, "Get pending order fail", http.StatusBadRequest)
		}

		if !pOrder.OtpExpiresAt.Valid || time.Now().After(pOrder.OtpExpiresAt.Time) {
			return utils.NewError(http.StatusBadRequest, "Otp expired")
		}

		if strings.TrimSpace(arg.Otp) != strings.TrimSpace(pOrder.Otp) {
			return utils.NewError(http.StatusBadRequest, "Otp invalid")
		}

		// 5. Tạo order thật
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
