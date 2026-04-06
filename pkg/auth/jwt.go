package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/cache"
	"time"
)

type JwtService struct {
	cache cache.RedisCacheService
}

func NewJwtService(cache cache.RedisCacheService) *JwtService {
	return &JwtService{
		cache: cache,
	}
}

const (
	AccessTokenTTl  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

var (
	jwtSecret  = []byte(utils.GetEnv("JWT_SECRET", "doan-vu-duy-hoc-lap-trinh-golang"))
	jwtEncrypt = []byte(utils.GetEnv("JWT_ENCRYPT", "12345678901234567890123456789012"))
)

type EncryptPayload struct {
	UserUuid pgtype.UUID `json:"user_uuid"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     string      `json:"role"`
}
type RefreshToken struct {
	UserUuid  pgtype.UUID `json:"user_uuid"`
	Token     string      `json:"token"`
	Email     string      `json:"email"`
	ExpiresAt time.Time   `json:"expires_at"`
	Revoked   bool        `json:"revoked"`
}

// ============= Access Token ==================
func (js *JwtService) GenerateToken(user sqlc.User) (string, error) {
	payload := EncryptPayload{
		UserUuid: user.UserUuid,
		Username: user.Name,
		Email:    user.Email,
		Role:     *user.Role,
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	data, err := EncryptAES(rawData, jwtSecret)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(AccessTokenTTl).Unix(),
		"iat":  time.Now().Unix(),
		"jti":  uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Đây là hàm phân tích token
func (js *JwtService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, nil, utils.WrapError(err, "Token is invalid!", utils.ErrorCodeUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.WrapError(err, "Claim is invalid!", utils.ErrorCodeUnauthorized)
	}
	return token, claims, nil
}

func (js *JwtService) EncryptAccessToken(tokenString string) (*EncryptPayload, error) {
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.WrapError(err, "Fail to parse Token!", utils.ErrorCodeUnauthorized)
	}
	data := claims["data"].(string)

	DecryptData, err := DecryptAES(data, jwtSecret)
	if err != nil {
		return nil, err
	}
	var payload EncryptPayload
	err = json.Unmarshal(DecryptData, &payload)
	return &payload, err
}

// ============= Refresh Token ==================
func (js *JwtService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenByte := make([]byte, 32)
	_, err := rand.Read(tokenByte)
	if err != nil {
		return RefreshToken{}, err
	}
	tokenString := base64.URLEncoding.EncodeToString(tokenByte)

	return RefreshToken{
		Token:     tokenString,
		UserUuid:  user.UserUuid,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil

}

func (js *JwtService) SaveRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JwtService) ValidateRefreshToken(tokenString string) (RefreshToken, error) {

	var refreshToken RefreshToken
	err := js.cache.Get("refresh_token:"+tokenString, &refreshToken)
	if err != nil || refreshToken.Revoked == true || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, err
	}
	return refreshToken, nil

}

func (js *JwtService) RevokedRefreshToken(token string) error {
	var refreshToken RefreshToken
	err := js.cache.Get("refresh_token:"+token, &refreshToken)
	if err != nil {
		return err
	}
	refreshToken.Revoked = true

	return js.cache.Set("refresh_token:"+token, refreshToken, time.Until(refreshToken.ExpiresAt))

}
