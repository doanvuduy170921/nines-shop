package utils

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	_ "golang.org/x/text/unicode/norm"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	}
}

func StringToInt32(paramStr string) int32 {
	str, err := strconv.Atoi(paramStr)
	if err != nil {
		log.Fatal(err)
	}
	return int32(str)
}

func StringToFloat64(paramStr string) float64 {
	str, err := strconv.Atoi(paramStr)
	if err != nil {
		log.Fatal(err)
	}
	return float64(str)
}

func StringToPgUuid(paramStr string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(paramStr); err != nil {
		return pgtype.UUID{}, err
	}
	uuid.Valid = true
	return uuid, nil

}

func PgTypeUuidToString(input pgtype.UUID) string {
	if !input.Valid {
		return ""
	}
	return input.String()
}

// cài go get github.com/gosimple/slug
func GenProductSlug(name string) string {
	if name == "" {
		return ""
	}
	return slug.Make(name)
}

// GenSKU tạo mã kho dựa trên tên và thuộc tính
// Ví dụ: Name: Apple, Attrs: {Color: Red, Size: XL} -> APPLE-RED-XL
func GenSKU(name string, attrs map[string]interface{}) string {
	parts := []string{strings.ToUpper(name[:3])} // Lấy 3 chữ đầu tên SP

	for _, v := range attrs {
		val := fmt.Sprintf("%v", v)
		if len(val) > 0 {
			parts = append(parts, strings.ToUpper(val))
		}
	}
	return strings.Join(parts, "-")
}

func Float64ToNumeric(val float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Format float64 thành string với đúng 2 chữ số thập phân.
	// Điều này biến 1.68872005e+08 thành "168872005.00"
	str := fmt.Sprintf("%.2f", val)

	// Sử dụng Scan để nạp chuỗi đã format vào pgtype.Numeric
	err := n.Scan(str)
	if err != nil {
		// Log lỗi nếu chuỗi không hợp lệ (trường hợp val là NaN hoặc Inf)
		log.Printf("❌ Numeric Scan Error: %v", err)
	}
	return n
}
