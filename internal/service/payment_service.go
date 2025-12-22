package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"math"
	"net/http"
	"net/url"
	"nineshop-be/internal/constant"
	"nineshop-be/internal/db"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/vnpay"
	"strconv"
)

type paymentService struct {
	repo repository.PaymentMethodRepository
}

func NewPaymentService(repo repository.PaymentMethodRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}

func (ps *paymentService) GetAllPayment(ctx *gin.Context) ([]sqlc.PaymentMethod, error) {
	c := ctx.Request.Context()
	return ps.repo.GetAllPayment(c)
}

func (ps *paymentService) HandleVnPaySuccess(c *gin.Context, orderRef string, query url.Values) error {
	ctx := c.Request.Context()
	pendingOrderID, err := strconv.Atoi(orderRef)
	if err != nil {
		return utils.WrapError(err, "Invalid order ref", http.StatusBadRequest)
	}

	tx, err := db.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return utils.WrapError(err, "Begin transaction failed", http.StatusInternalServerError)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	qTx := db.DB.WithTx(tx)

	// 1. Get pending order
	pOrder, err := qTx.GetPOrderById(ctx, int32(pendingOrderID))
	if err != nil {
		return utils.WrapError(err, "Get pending order fail", http.StatusBadRequest)
	}

	// 2. Create order
	order, err := qTx.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID:          pOrder.UserID,
		Name:            pOrder.Name,
		Email:           pOrder.Email,
		Phone:           pOrder.Phone,
		Address:         pOrder.Address,
		PaymentMethodID: pOrder.PaymentMethodID,
		ShippingPrice:   pOrder.ShippingPrice,
		Subtotal:        pOrder.Subtotal,
		Tax:             pOrder.Tax,
		TotalAmount:     pOrder.TotalAmount,
		Status:          constant.OrderStatusProcessing,
		AmountItem:      pOrder.AmountItem,
	})
	if err != nil {
		return utils.WrapError(err, "Create order fail", http.StatusInternalServerError)
	}
	txnRef := query.Get("vnp_TxnRef")
	// 3. Update payment info
	err = qTx.UpdateOrderPayment(ctx, sqlc.UpdateOrderPaymentParams{
		ID:            order.ID,
		PaymentStatus: constant.PaymentStatusPaid,
		TransactionID: &txnRef,
	})
	if err != nil {
		return utils.WrapError(err, "Update payment status fail", http.StatusInternalServerError)
	}

	// 4. Status history
	note := constant.OrderStatusNotes[constant.OrderStatusProcessing]
	_, err = qTx.CreateOrderStatusHistory(ctx, sqlc.CreateOrderStatusHistoryParams{
		OrderID: order.ID,
		Status:  constant.OrderStatusProcessing,
		Note:    &note,
	})
	if err != nil {
		return utils.WrapError(err, "Create order status history fail", http.StatusInternalServerError)
	}

	// 5. Copy items
	items, err := qTx.GetByPOrderItemId(ctx, pOrder.ID)
	if err != nil {
		return utils.WrapError(err, "Get pending order items fail", http.StatusInternalServerError)
	}

	for _, item := range items {
		product, err := qTx.GetProductById(ctx, item.ProductID)
		if err != nil {
			return err
		}

		_, err = qTx.AddOrderItem(ctx, sqlc.AddOrderItemParams{
			OrderID:          order.ID,
			ProductID:        item.ProductID,
			Quantity:         item.Quantity,
			Price:            item.Price,
			ProductName:      product.Name,
			ProductThumbnail: &product.Thumbnail,
		})
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	return nil
}

func (ps *paymentService) ProcessVNPayPayment(ctx *gin.Context, orderID int64, amount int64, transactionNo string, bankCode string) (sqlc.Order, error) {
	c := ctx.Request.Context()

	tx, err := db.DBPool.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return sqlc.Order{}, utils.NewError(http.StatusInternalServerError, "Begin transaction failed")
	}
	qTx := db.DB.WithTx(tx)

	defer func() {
		if err != nil {
			_ = tx.Rollback(c)
		}
	}()

	// 2. Lấy pending order theo ID
	pOrder, err := qTx.GetPOrderById(c, int32(orderID))
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Pending order not found", http.StatusBadRequest)
	}

	// 3. Kiểm tra pending order có phải VNPay không
	if pOrder.PaymentMethodID != constant.PaymentMethodVnPay {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "This order is not VNPay payment")
	}

	// 5. Kiểm tra số tiền có khớp không
	vnpayAmountXu := amount         // amount từ vnpay.VerifyAndParseIPN đã chia 100 rồi
	vnpayAmountVND := vnpayAmountXu // Đây là VND

	floatVal, _ := pOrder.TotalAmount.Float64Value()
	expectedAmountUSD := floatVal.Float64
	expectedAmountVND := int64(math.Round(expectedAmountUSD * vnpay.USD_TO_VND_RATE))

	log.Printf("🔍 Amount check: VNPay=%d VND, Expected=%d VND", vnpayAmountVND, expectedAmountVND)

	// ✅ Cho phép chênh lệch ±1 VND do làm tròn
	diff := vnpayAmountVND - expectedAmountVND
	if diff < -1 || diff > 1 {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest,
			fmt.Sprintf("Amount mismatch: VNPay=%d VND, Expected=%d VND (diff=%d)",
				vnpayAmountVND, expectedAmountVND, diff))
	}

	log.Printf("✅ Amount matched (diff=%d VNĐ)", diff)

	// 6. Tạo order thật
	paymentStatus := constant.PaymentStatusPaid
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
		PaymentStatus:   paymentStatus,
		TransactionID:   &transactionNo,
		AmountItem:      pOrder.AmountItem,
	})
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Create order fail", http.StatusInternalServerError)
	}

	// 7. Tạo order status history
	statusHistory := constant.OrderStatusPending
	note := constant.OrderStatusNotes[statusHistory]
	_, err = qTx.CreateOrderStatusHistory(c, sqlc.CreateOrderStatusHistoryParams{
		OrderID: order.ID,
		Status:  statusHistory,
		Note:    &note,
	})
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Create order status history fail", http.StatusInternalServerError)
	}

	// 8. Lấy list pending order items
	pOrderItems, err := qTx.GetByPOrderItemId(c, pOrder.ID)
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Get pending order items fail", http.StatusInternalServerError)
	}

	// 9. Copy từng item sang order_items
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

	// 10. Commit transaction
	if err := tx.Commit(c); err != nil {
		return sqlc.Order{}, utils.NewError(http.StatusInternalServerError, "Commit transaction failed")
	}

	log.Printf("✅✅✅ VNPay: Order created successfully - ID=%d, TransactionNo=%s", order.ID, transactionNo)

	return order, nil
}
