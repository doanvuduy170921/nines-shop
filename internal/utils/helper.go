package utils

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	_ "golang.org/x/text/unicode/norm"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
		log.Fatal("Error loading .env file")
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

func GenProductSku(name string) string {
	if name == "" {
		return ""
	}
	nomalized := slug.MakeLang(name, "en")

	prefix := strings.ToUpper(strings.ReplaceAll(nomalized, " ", ""))

	dateAt := time.Now().Format("060201")
	randPart := rand.Intn(9999)

	return fmt.Sprintf("%s-%s-%04d", prefix, dateAt, randPart)
}
