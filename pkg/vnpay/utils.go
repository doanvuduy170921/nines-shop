package vnpay

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func CreatePaymentURL(orderID int64, amount int64, ipAddr string) (string, error) {
	cfg := LoadConfig()

	params := map[string]string{
		"vnp_Version":    Version,
		"vnp_Command":    Command,
		"vnp_TmnCode":    TmnCode,
		"vnp_Amount":     strconv.FormatInt(amount, 10), // nhân 100
		"vnp_CreateDate": time.Now().Format("20060102150405"),
		"vnp_CurrCode":   CurrCode,
		"vnp_IpAddr":     ipAddr,
		"vnp_Locale":     Locale,
		"vnp_OrderInfo":  "Thanh toan don hang #" + strconv.FormatInt(orderID, 10),
		"vnp_OrderType":  "other",
		"vnp_ReturnUrl":  cfg.ReturnURL,
		"vnp_TxnRef":     TxRefPrefix + strconv.FormatInt(orderID, 10),
	}

	// Tạo chuỗi query + hash
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var query strings.Builder
	for i, k := range keys {
		if i > 0 {
			query.WriteByte('&')
		}
		query.WriteString(url.QueryEscape(k) + "=" + url.QueryEscape(params[k]))
	}

	hash := hmacSHA512(query.String(), HashSecret)
	return PaymentURL + "?" + query.String() + "&vnp_SecureHash=" + hash, nil
}

func VerifySecureHash(query url.Values) bool {
	receivedHash := query.Get("vnp_SecureHash")
	if receivedHash == "" {
		return false
	}
	query.Del("vnp_SecureHash")
	query.Del("vnp_SecureHashType")

	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data strings.Builder
	for i, k := range keys {
		if i > 0 {
			data.WriteByte('&')
		}
		data.WriteString(url.QueryEscape(k) + "=" + url.QueryEscape(query.Get(k)))
	}

	expected := hmacSHA512(data.String(), HashSecret)
	return hmac.Equal([]byte(receivedHash), []byte(expected))
}

func hmacSHA512(data, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
