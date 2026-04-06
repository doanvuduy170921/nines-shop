package service

import (
	"context"
	"log"
	"net/http"
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
	store *db.Store
	repo  repository.PendingOrderRepository
	ui    UserInject
	pot   repository.PendingOrderItemRepository
	email email.EmailService
	or    repository.OrderRepository
	oi    repository.OrderItemRepository
	pr    repository.ProductRepository
	osh   repository.OrderStatusHistoryRepository
}

type UserInject interface {
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
}

func NewPendingOrderService(store *db.Store, repo repository.PendingOrderRepository, ui UserInject, pot repository.PendingOrderItemRepository, email email.EmailService, or repository.OrderRepository, oi repository.OrderItemRepository, pr repository.ProductRepository, osh repository.OrderStatusHistoryRepository) PendingOrderService {
	return &pendingOrderService{
		store: store,
		repo:  repo,
		ui:    ui,
		pot:   pot,
		email: email,
		or:    or,
		oi:    oi,
		pr:    pr,
		osh:   osh,
	}
}

func (ns *pendingOrderService) Create(ctx *gin.Context, arg dto.CreatePendingOrderDto) (res dto.CreatePendingOrderResponse, err error) {
	c := ctx.Request.Context()

	userUuid, exists := ctx.Get("user_uuid")
	if !exists {
		return dto.CreatePendingOrderResponse{}, utils.NewError(400, "user uuid not found in request context")
	}
	userUuidStr := userUuid.(string)
	Uuid, err := utils.StringToPgUuid(userUuidStr)
	if err != nil {
		return dto.CreatePendingOrderResponse{}, utils.NewError(400, "convert user uuid to uuid fail")
	}

	if len(arg.Items) == 0 {
		return dto.CreatePendingOrderResponse{}, utils.NewError(400, "items cannot be empty")
	}

	var status = constant.OrderStatusPending

	var (
		user           sqlc.User
		otp            string
		otpExpiresAt   pgtype.Timestamp
		shouldSendOTP  bool
		pOrder         sqlc.PendingOrder
		paymentURL     string
		totalAmount    float64
		vnpayRequested bool
	)

	err = ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		var errTx error

		user, errTx = q.GetByUuid(c, Uuid)
		if errTx != nil {
			return utils.NewError(http.StatusBadRequest, "get user by uuid fail")
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

		// đếm số sản phẩm
		count := int32(len(arg.Items))

		log.Printf("DEBUG PENDING ORDER: Subtotal=%v, Total=%v, Tax=%v", arg.Subtotal, arg.Total, arg.Tax)

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

		if vnpayRequested {
			totalAmountNumeric, _ := pOrder.TotalAmount.Float64Value()
			totalAmount = totalAmountNumeric.Float64
		}

		return nil
	})

	if err != nil {
		return dto.CreatePendingOrderResponse{}, err
	}

	// Sau khi transaction thành công:
	// 1. Nếu VNPAY: tạo URL thanh toán
	if vnpayRequested {
		vnPayService := vnpay.NewService()
		paymentURL, err = vnPayService.CreatePaymentURL(vnpay.OrderInfo{
			ID:          int64(pOrder.ID),
			TotalAmount: totalAmount,
		}, ctx.ClientIP())
		if err != nil {
			log.Printf("ERROR creating VNPAY URL: %v", err)
			return dto.CreatePendingOrderResponse{}, utils.WrapError(err, "Create VNPAY URL fail", http.StatusInternalServerError)
		}
	}

	// 2. Nếu COD: gửi email OTP đồng bộ, lấy email từ DB (user.Email)
	if shouldSendOTP {
		// Tôn trọng context: nếu request đã bị hủy thì không gửi nữa
		select {
		case <-c.Done():
			log.Printf("skip sending OTP email, context canceled for user %d", user.ID)
		default:
			if sendErr := ns.email.SendOTPEmail(user.Email, otp); sendErr != nil {
				log.Printf("send OTP email fail for user %d (%s): %v", user.ID, user.Email, sendErr)
				return dto.CreatePendingOrderResponse{}, utils.WrapError(sendErr, "Send OTP email fail", http.StatusInternalServerError)
			}
		}
	}

	return dto.CreatePendingOrderResponse{
		PendingOrder: pOrder,
		VnPayURL:     paymentURL,
	}, nil
}

func (ns *pendingOrderService) ValidateOTP(ctx *gin.Context, arg dto.ValidateOTPParams) (sqlc.Order, error) {
	c := ctx.Request.Context()

	var order sqlc.Order

	err := ns.store.ExecTx(c, func(q *sqlc.Queries) error {
		// 2. Lấy pending order theo pId
		pOrder, errTx := q.GetPOrderById(c, int32(arg.POrderId))
		if errTx != nil {
			return utils.WrapError(errTx, "Get pending order fail", http.StatusBadRequest)
		}

		// 3. Kiểm tra mã OTP hết hạn hay chưa
		if !pOrder.OtpExpiresAt.Valid || time.Now().After(pOrder.OtpExpiresAt.Time) {
			return utils.NewError(http.StatusBadRequest, "Otp expired")
		}

		// 4. Kiểm tra OTP có hợp lệ không
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
		// 6. Tạo order status history đầu tiên
		_, errTx = q.CreateOrderStatusHistory(c, sqlc.CreateOrderStatusHistoryParams{
			OrderID: order.ID,
			Status:  statusHistory,
			Note:    &note,
		})
		if errTx != nil {
			return utils.WrapError(errTx, "Create order status history fail", http.StatusInternalServerError)
		}

		// 7. Lấy list pending order item
		pOrderItems, errTx := q.GetByPOrderItemId(c, pOrder.ID)
		if errTx != nil {
			return utils.WrapError(errTx, "Get pending order items fail", http.StatusInternalServerError)
		}

		// 8. Copy từng item sang order_items
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
