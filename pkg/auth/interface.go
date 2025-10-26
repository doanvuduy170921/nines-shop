package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"nineshop-be/internal/db/sqlc"
)

type TokenService interface {
	GenerateToken(user sqlc.User) (string, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	EncryptAccessToken(tokenString string) (*EncryptPayload, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	SaveRefreshToken(token RefreshToken) error
	ValidateRefreshToken(tokenString string) (RefreshToken, error)
	RevokedRefreshToken(token string) error
}
