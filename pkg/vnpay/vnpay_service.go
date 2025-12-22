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
	// Quy đổi USD sang VND
	amountVND := order.TotalAmount * USD_TO_VND_RATE

	// Làm tròn và nhân 100 để ra đơn vị xu
	amountInXu := int64(math.Round(amountVND * 100))

	// Debug log
	log.Printf("💰 CreatePaymentURL: $%.2f USD → %.0f VND → %d xu",
		order.TotalAmount, amountVND, amountInXu)

	// Kiểm tra giới hạn VNPay (5,000 - 999,999,999 VND)
	if amountVND < 5000 {
		return "", fmt.Errorf("Amount too small: %.0f VND (min 5,000 VND)", amountVND)
	}
	if amountVND >= 1000000000 {
		return "", fmt.Errorf("Amount too large: %.0f VND (max 999,999,999 VND)", amountVND)
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
