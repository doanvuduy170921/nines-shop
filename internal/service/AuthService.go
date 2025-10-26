package service

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/auth"
	"nineshop-be/pkg/cache"
	"strings"
	"time"
)

type authService struct {
	repo         repository.UserRepository
	tokenService auth.TokenService
	cache        cache.RedisCacheService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cache cache.RedisCacheService) AuthService {
	return &authService{
		repo:         repo,
		tokenService: tokenService,
		cache:        cache,
	}
}

func (as *authService) Login(c *gin.Context, email, password string) (string, string, int, string, error) {
	context := c.Request.Context()
	user, err := as.repo.FindByEmail(context, email)
	if err != nil {
		return "", "", 0, "", utils.WrapError(err, "Username or password is invalid! (1)", utils.ErrorCodeUnauthorized)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", 0, "", utils.WrapError(err, "Username or password is invalid! (2)", utils.ErrorCodeUnauthorized)
	}
	accessToken, err := as.tokenService.GenerateToken(user)
	if err != nil {
		return "", "", 0, "", utils.WrapError(err, "Fail to generate access token", utils.ErrorCodeInternalError)
	}
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, "", utils.WrapError(err, "Fail to generate refresh token", utils.ErrorCodeInternalError)
	}

	err = as.tokenService.SaveRefreshToken(refreshToken)
	if err != nil {
		log.Println("⛔ Error saving refresh token", err)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTl), *user.Role, nil

}

func (as *authService) RefreshToken(c *gin.Context, tokenString string) (string, string, int, error) {
	context := c.Request.Context()
	// Kiểm tra refreshToken , trả thông tin cho user email
	token, err := as.tokenService.ValidateRefreshToken(tokenString)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "refresh token is invalid", utils.ErrorCodeForbidden)
	}
	// lấy thông tin của user
	user, err := as.repo.FindByEmail(context, token.Email)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Fail to find by email ", utils.ErrorCodeNotFound)
	}

	// tạo access token mới
	accessToken, err := as.tokenService.GenerateToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Fail to generate access token", utils.ErrorCodeInternalError)
	}
	// tạo refreshToken mới
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Fail to generate refresh token", utils.ErrorCodeInternalError)
	}

	// vô hiệu hóa refreshToken cũ
	err = as.tokenService.RevokedRefreshToken(tokenString)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Fail to revoke refresh token", utils.ErrorCodeForbidden)
	}

	//Lưu refreshToken mới
	err = as.tokenService.SaveRefreshToken(refreshToken)
	if err != nil {
		log.Println("⛔ Error saving refresh token", err)
	}

	return accessToken, refreshToken.Token, int(auth.RefreshTokenTTL), nil

}

func (as *authService) Logout(c *gin.Context, refreshTokenString string) error {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return nil
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	_, claims, err := as.tokenService.ParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return nil
	}
	if jti, ok := claims["jti"].(string); ok {
		expUnix := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		err = as.cache.Set(key, "revoked", ttl)
	}
	// validate refreshToken
	_, err = as.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return utils.WrapError(err, "Fail to validate refresh token", utils.ErrorCodeForbidden)
	}
	// revoked refreshToken
	err = as.tokenService.RevokedRefreshToken(refreshTokenString)
	if err != nil {
		return utils.WrapError(err, "Fail to revoke refresh token", utils.ErrorCodeForbidden)
	}
	return nil
}
