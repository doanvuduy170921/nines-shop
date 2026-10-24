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
		product, err := qTx.GetVariantById(ctx, item.VariantID)
		if err != nil {
			return err
		}

		_, err = qTx.AddOrderItem(ctx, sqlc.AddOrderItemParams{
			OrderID:          order.ID,
			VariantID:        item.VariantID,
			Quantity:         item.Quantity,
			Price:            item.Price,
			ProductName:      product.ProductName,
			ProductThumbnail: &product.ProductThumbnail,
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

	pOrder, err := qTx.GetPOrderById(c, int32(orderID))
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Pending order not found", http.StatusBadRequest)
	}

	if pOrder.PaymentMethodID != constant.PaymentMethodVnPay {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest, "This order is not VNPay payment")
	}

	vnpayAmountVND := amount

	floatVal, _ := pOrder.TotalAmount.Float64Value()
	expectedAmountVNDFloat := floatVal.Float64

	expectedAmountVND := int64(math.Round(expectedAmountVNDFloat))

	log.Printf("🔍 Amount check: VNPay=%d VND, Expected=%d VND", vnpayAmountVND, expectedAmountVND)

	diff := vnpayAmountVND - expectedAmountVND
	if diff < -1 || diff > 1 {
		return sqlc.Order{}, utils.NewError(http.StatusBadRequest,
			fmt.Sprintf("Amount mismatch: VNPay=%d VND, Expected=%d VND (diff=%d)",
				vnpayAmountVND, expectedAmountVND, diff))
	}

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

	pOrderItems, err := qTx.GetByPOrderItemId(c, pOrder.ID)
	if err != nil {
		return sqlc.Order{}, utils.WrapError(err, "Get pending order items fail", http.StatusInternalServerError)
	}

	for _, item := range pOrderItems {
		variant, err := qTx.GetVariantForUpdate(c, item.VariantID)
		if err != nil {
			return sqlc.Order{}, utils.WrapError(err, fmt.Sprintf("Variant %d does not exist", item.VariantID), http.StatusBadRequest)
		}

		if !*variant.IsActive {
			return sqlc.Order{}, utils.NewError(http.StatusBadRequest, fmt.Sprintf("Product %s is no longer available for sale", variant.ProductName))
		}

		res, err := qTx.DecreaseStock(c, sqlc.DecreaseStockParams{
			VariantID: item.VariantID,
			Stock:     item.Quantity,
		})
		if err != nil {
			return sqlc.Order{}, utils.WrapError(err, "Inventory update error", http.StatusInternalServerError)
		}

		rowsAffected := res.RowsAffected()
		if rowsAffected == 0 {
			return sqlc.Order{}, utils.NewError(http.StatusBadRequest,
				fmt.Sprintf("Product %s has insufficient stock (Remaining: %d)", variant.ProductName, variant.StockQuantity))
		}

		_, err = qTx.AddOrderItem(c, sqlc.AddOrderItemParams{
			OrderID:          order.ID,
			VariantID:        item.VariantID,
			Quantity:         item.Quantity,
			Price:            item.Price,
			ProductName:      variant.ProductName,
			ProductThumbnail: &variant.ProductThumbnail,
		})
		if err != nil {
			return sqlc.Order{}, utils.WrapError(err, "Add order item fail", http.StatusInternalServerError)
		}
	}
	if err := tx.Commit(c); err != nil {
		return sqlc.Order{}, utils.NewError(http.StatusInternalServerError, "Commit transaction failed")
	}

	log.Printf("✅✅✅ VNPay: Order created successfully - ID=%d, TransactionNo=%s", order.ID, transactionNo)

	return order, nil
}
