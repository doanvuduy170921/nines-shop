package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/pkg/auth"
	"nineshop-be/pkg/cache"
	"strings"
)

func AuthMiddleware(tokenService auth.TokenService, cache cache.RedisCacheService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized header missing or invalid",
			})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := tokenService.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
		if jti, ok := claims["jti"].(string); ok {
			key := "blacklist:" + jti
			if exists, err := cache.Exists(key); err == nil && exists {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token exists",
				})
				return
			}
		}

		payload, err := tokenService.EncryptAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
		c.Set("user_name", payload.Username)
		c.Set("email", payload.Email)
		c.Set("role", payload.Role)
		c.Set("user_uuid", payload.UserUuid)

		c.Next()
	}
}
