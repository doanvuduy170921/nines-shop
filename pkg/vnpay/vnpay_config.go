package vnpay

import "os"

type Config struct {
	ReturnURL string // frontend nhận kết quả
	IPNURL    string // backend nhận callback (không cần dùng trong utils)
}

func LoadConfig() Config {
	return Config{
		ReturnURL: getEnv("VNPAY_RETURN_URL", "http://localhost:8080/api/v1/payment/vnpay/return"),
		IPNURL:    getEnv("VNPAY_IPN_URL", "http://localhost:8080/api/v1/payment/vnpay/ipn"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
