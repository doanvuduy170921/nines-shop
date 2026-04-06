package vnpay

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"
	"math/big"
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
		"vnp_Amount":     strconv.FormatInt(amount, 10),
		"vnp_CreateDate": time.Now().Format("20060102150405"),
		"vnp_CurrCode":   CurrCode,
		"vnp_IpAddr":     ipAddr,
		"vnp_Locale":     Locale,
		"vnp_OrderInfo":  "Thanh toan don hang " + strconv.FormatInt(orderID, 10),
		"vnp_OrderType":  "other",
		"vnp_ReturnUrl":  cfg.ReturnURL,
		"vnp_TxnRef":     TxRefPrefix + strconv.FormatInt(orderID, 10),
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ✅ FIX: Tạo hashData KHÔNG encode - đúng chuẩn VNPAY
	var hashData strings.Builder
	var queryBuilder strings.Builder

	for i, k := range keys {
		if i > 0 {
			hashData.WriteByte('&')
			queryBuilder.WriteByte('&')
		}
		// Hash: raw value, không encode
		hashData.WriteString(k + "=" + params[k])
		// URL: encode value, space → %20 (không phải +)
		encodedVal := strings.ReplaceAll(url.QueryEscape(params[k]), "+", "%20")
		queryBuilder.WriteString(k + "=" + encodedVal)
	}

	hash := hmacSHA512(hashData.String(), HashSecret)

	// 👇 DEBUG LOG
	log.Printf("🔐 Hash input:\n%s", hashData.String())
	log.Printf("🔐 Hash output (%d chars): %s", len(hash), hash)

	queryBuilder.WriteString("&vnp_SecureHash=" + hash)

	return PaymentURL + "?" + queryBuilder.String(), nil
}

func VerifySecureHash(query url.Values) bool {
	receivedHash := query.Get("vnp_SecureHash")
	if receivedHash == "" {
		return false
	}

	// Clone để không mutate original
	cloned := url.Values{}
	for k, v := range query {
		cloned[k] = v
	}
	cloned.Del("vnp_SecureHash")
	cloned.Del("vnp_SecureHashType")

	var keys []string
	for k := range cloned {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ✅ FIX: Hash trên raw string, KHÔNG encode - đúng chuẩn VNPAY
	var data strings.Builder
	for i, k := range keys {
		if i > 0 {
			data.WriteByte('&')
		}
		data.WriteString(k + "=" + cloned.Get(k))
	}

	expected := hmacSHA512(data.String(), HashSecret)
	return hmac.Equal([]byte(strings.ToLower(receivedHash)), []byte(expected))
}

func hmacSHA512(data, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(data))
	result := mac.Sum(nil)
	// Đảm bảo luôn 128 ký tự
	hash := fmt.Sprintf("%064x%064x",
		new(big.Int).SetBytes(result[:32]),
		new(big.Int).SetBytes(result[32:]))
	log.Printf("🔐 Hash length: %d", len(hash))
	return hash
}
