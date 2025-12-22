package vnpay

import "net/url"

type Service interface {
	CreatePaymentURL(order OrderInfo, ipAddr string) (string, error)
	VerifyAndParseIPN(query url.Values) (orderID int64, amount int64, success bool, err error)
}
