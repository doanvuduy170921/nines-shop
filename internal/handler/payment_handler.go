package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/vnpay"
	"strings"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (ph *PaymentHandler) GetAllPayment(ctx *gin.Context) {

	payments, err := ph.service.GetAllPayment(ctx)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, payments, "Get All Payment")
}

func (ph *PaymentHandler) CallBackFromVnPay(ctx *gin.Context) {
	query := ctx.Request.URL.Query()

	if !vnpay.VerifySecureHash(query) {
		ctx.Redirect(302, "/payment/failed")
		return
	}

	if query.Get("vnp_ResponseCode") != "00" {
		ctx.Redirect(302, "/payment/failed")
		return
	}

	txnRef := query.Get("vnp_TxnRef")
	orderRef := strings.TrimPrefix(txnRef, vnpay.TxRefPrefix)
	err := ph.service.HandleVnPaySuccess(ctx, orderRef, query)
	if err != nil {
		ctx.Redirect(302, "/payment/failed")
		return
	}

	ctx.Redirect(302, "/payment/success?order_id="+orderRef)
}

func (ph *PaymentHandler) Redirect(ctx *gin.Context) {

	// 👇 THÊM 2 DÒNG NÀY ĐỂ DEBUG
	log.Printf("🔍 RAW URL: %s", ctx.Request.URL.RawQuery)
	log.Printf("🔍 Parsed Query: %v", ctx.Request.URL.Query())

	query := ctx.Request.URL.Query()

	// 1. Verify và parse
	vnpayService := vnpay.NewService()
	orderID, amount, success, err := vnpayService.VerifyAndParseIPN(query)

	if err != nil || !success {
		log.Printf("VNPay return: Invalid - %v", err)
		ctx.Redirect(http.StatusFound, "http://localhost:4200/payment-result?status=failed")
		return
	}

	// 2. TẠO ORDER THẬT LUÔN Ở ĐÂY
	transactionNo := query.Get("vnp_TransactionNo")
	bankCode := query.Get("vnp_BankCode")

	order, err := ph.service.ProcessVNPayPayment(ctx, orderID, amount, transactionNo, bankCode)
	if err != nil {
		log.Printf("Create order error: %v", err)
		// Nếu đã tạo rồi, vẫn cho qua
		if strings.Contains(err.Error(), "already confirmed") {
			ctx.Redirect(http.StatusFound, fmt.Sprintf("http://localhost:4200/payment-result?status=success&order_id=%d", orderID))
			return
		}

		ctx.Redirect(http.StatusFound, "http://localhost:4200/payment-result?status=failed")
		return
	}

	// 3. Success
	log.Printf("✅ Order created: ID=%d", order.ID)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("http://localhost:4200/payment-result?status=success&order_id=%d", order.ID))
}
