package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"nineshop-be/internal/constant"
	"nineshop-be/internal/db"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/email"
	"strings"
	"time"
)

type pendingOrderService struct {
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

func NewPendingOrderService(repo repository.PendingOrderRepository, ui UserInject, pot repository.PendingOrderItemRepository, email email.EmailService, or repository.OrderRepository, oi repository.OrderItemRepository, pr repository.ProductRepository, osh repository.OrderStatusHistoryRepository) PendingOrderService {
	return &pendingOrderService{
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

func (ns *pendingOrderService) Create(ctx *gin.Context, arg dto.CreatePendingOrderDto) (sqlc.PendingOrder, error) {
	c := ctx.Request.Context()
	userUuid, exists := ctx.Get("user_uuid")
	if !exists {
		return sqlc.PendingOrder{}, utils.NewError(400, "user uuid not found in request context")
	}
	userUuidStr := userUuid.(string)
	Uuid, err := utils.StringToPgUuid(userUuidStr)
	if err != nil {
		return sqlc.PendingOrder{}, utils.NewError(400, "convert user uuid to uuid fail")
	}
	user, err := ns.ui.GetByUuid(c, Uuid)
	if err != nil {
		return sqlc.PendingOrder{}, utils.NewError(400, "get user by uuid fail")
	}

	if len(arg.Items) == 0 {
		return sqlc.PendingOrder{}, utils.NewError(400, "items cannot be empty")
	}

	otp := utils.GenOTP()

	var otpExpiresAt pgtype.Timestamp

	otpExpiresAt = pgtype.Timestamp{
		Time:  time.Now().Add(5 * time.Minute),
		Valid: true,
	}

	var status = "pending"

	// send otp by email
	go func(email, otp string) {
		if err := ns.email.SendOTPEmail(email, otp); err != nil {
			log.Printf("send email fail: %v", err)
		}
	}(arg.Email, otp) // gọi hàm bằng cách truyền tham số, ở trên định ngĩa func

	pOrder, err := ns.repo.Create(c, sqlc.CreatePendingOrderParams{
		UserID:          user.ID,
		Name:            arg.Name,
		Email:           arg.Email,
		Phone:           arg.Phone,
		PaymentMethodID: int32(arg.PaymentId),
		Address:         arg.Address,
		Otp:             otp,
		OtpExpiresAt:    otpExpiresAt,
		Subtotal:        utils.Float64ToPgTypeNumeric(arg.Subtotal),
		TotalAmount:     utils.Float64ToPgTypeNumeric(arg.Total),
		Tax:             utils.Float64ToPgTypeNumeric(arg.Tax),
		ShippingPrice:   utils.Float64ToPgTypeNumeric(arg.Shipping),
		Status:          &status,
	})
	if err != nil {
		log.Printf("ERROR creating pending order: %v", err)
		return sqlc.PendingOrder{}, utils.WrapError(err, "Create Pending order fail", 400)
	}

	for _, item := range arg.Items {
		_, err := ns.pot.Create(c, sqlc.CreatePendingOrderItemParams{
			PendingOrderID: pOrder.ID,
			ProductID:      int32(item.ProductId),
			Quantity:       int32(item.Quantity),
			Price:          utils.Float64ToPgTypeNumeric(item.Price),
		})
		if err != nil {
			return sqlc.PendingOrder{}, utils.WrapError(err, "Create item fail", 400)
		}
	}

	return pOrder, nil
}

func (ns *pendingOrderService) ValidateOTP(ctx *gin.Context, arg dto.ValidateOTPParams) (sqlc.Order, error) {
	c := ctx.Request.Context()

	// 1. Bắt đầu transaction
	tx, err := db.DBPool.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return sqlc.Order{}, utils.NewError(http.StatusInternalServerError, "Begin transaction failed")
	}
	qTx := db.DB.WithTx(tx)

	// Nếu err != nil ở bất kỳ chỗ nào -> rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback(c)
		}
	}()

	// 2. Lấy pending order theo pId
	pOrder, err := qTx.GetPOrderById(c, int32(arg.POrderId))
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Get pending order fail", http.StatusBadRequest)
	}

	// 3. Kiểm tra mã OTP hết hạn hay chưa
	if !pOrder.OtpExpiresAt.Valid || time.Now().After(pOrder.OtpExpiresAt.Time) {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "Otp expired")
	}

	// 4. Kiểm tra OTP có hợp lệ không
	if strings.TrimSpace(arg.Otp) != strings.TrimSpace(pOrder.Otp) {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "Otp invalid")
	}

	// 5. Tạo order thật
	order, err := qTx.CreateOrder(c, sqlc.CreateOrderParams{
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
	})
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Create order fail", http.StatusInternalServerError)
	}

	// 6. Tạo order status history đầu tiên
	_, err = qTx.CreateOrderStatusHistory(c, sqlc.CreateOrderStatusHistoryParams{
		OrderID: order.ID,
		Status:  constant.OrderStatusPending,
	})
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Create order status history fail", http.StatusInternalServerError)
	}

	// 7. Lấy list pending order item
	pOrderItems, err := qTx.GetByPOrderItemId(c, pOrder.ID)
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Get pending order items fail", http.StatusInternalServerError)
	}

	// 8. Copy từng item sang order_items
	for _, item := range pOrderItems {
		product, err := qTx.GetProductById(c, item.ProductID)
		if err != nil {
			return sqlc.Order{}, utils.WrapError(err, "Get product from pending order item fail", http.StatusBadRequest)
		}

		_, err = qTx.AddOrderItem(c, sqlc.AddOrderItemParams{
			OrderID:          order.ID,
			ProductID:        item.ProductID,
			Quantity:         item.Quantity,
			Price:            item.Price,
			ProductName:      product.Name,
			ProductThumbnail: &product.Thumbnail,
		})
		if err != nil {
			return sqlc.Order{}, utils.WrapError(err, "Add order item fail", http.StatusInternalServerError)
		}
	}

	// 9. Commit transaction
	if err := tx.Commit(c); err != nil {
		return sqlc.Order{}, utils.NewError(http.StatusInternalServerError, "Commit transaction failed")
	}

	return order, nil
}
