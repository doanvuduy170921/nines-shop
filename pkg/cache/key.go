package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const otpPrefix = "otp"

func HashEmail(email string) string {
	h := sha256.New()
	h.Write([]byte(strings.ToLower(email)))
	return hex.EncodeToString(h.Sum(nil))
}

func OTPCreate(email string) string {
	return fmt.Sprintf("otp:create:%s", HashEmail(email))
}
