package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func XApiKeyMiddlerWare() gin.HandlerFunc {
	acceptApiKey := os.Getenv("X_API_KEY")
	if acceptApiKey == "" {
		acceptApiKey = "secret key"
	}
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X_API_KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Missing X-API-KEY header",
			})
			return
		}

		if apiKey != acceptApiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid API key",
			})
			return
		}
		c.Next()
	}

}
