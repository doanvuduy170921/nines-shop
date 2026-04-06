package vnpay

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type OrderInfo struct {
	ID          int64
	TotalAmount float64 // vì pgtype.Numeric có method .Float64
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) CreatePaymentURL(order OrderInfo, ipAddr string) (string, error) {
	// ✅ FIX: Convert IPv6 localhost → IPv4
	if ipAddr == "::1" || ipAddr == "" {
		ipAddr = "127.0.0.1"
	}
	
	// 1. Số tiền bây giờ đã là VND trực tiếp từ Frontend/DB
	amountVND := order.TotalAmount

	// 2. VNPAY yêu cầu số tiền nhân 100 để ra đơn vị "xu" (cents)
	// Lưu ý: Dùng math.Round để đảm bảo không bị sai số lẻ
	amountInXu := int64(math.Round(amountVND * 100))

	// Debug log - Cập nhật lại log cho chính xác đơn vị
	log.Printf("💰 CreatePaymentURL: %.0f VND → %d xu", amountVND, amountInXu)

	// 3. Kiểm tra giới hạn VNPay (5,000 - 999,999,999 VND)
	if amountVND < 5000 {
		return "", fmt.Errorf("Số tiền quá nhỏ: %.0f VND (tối thiểu 5,000 VND)", amountVND)
	}
	// Giới hạn 1 tỷ đồng (VNPAY thường giới hạn giao dịch dưới 1 tỷ)
	if amountVND >= 1000000000 {
		return "", fmt.Errorf("Số tiền quá lớn: %.0f VND (tối đa 999,999,999 VND)", amountVND)
	}

	return CreatePaymentURL(order.ID, amountInXu, ipAddr)
}

func (s *service) VerifyAndParseIPN(query url.Values) (orderID int64, amount int64, success bool, err error) {
	if !VerifySecureHash(query) {
		return 0, 0, false, nil
	}

	if query.Get("vnp_ResponseCode") != "00" {
		return 0, 0, false, nil
	}

	txnRef := query.Get("vnp_TxnRef")
	orderIDStr := strings.TrimPrefix(txnRef, TxRefPrefix)
	orderID, err = strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return 0, 0, false, err
	}

	amt, _ := strconv.ParseInt(query.Get("vnp_Amount"), 10, 64)
	amount = amt / 100

	return orderID, amount, true, nil
}
