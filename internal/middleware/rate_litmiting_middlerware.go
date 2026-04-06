package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"nineshop-be/pkg/cache"
	"time"
)

func RateLimitMiddleware(rdb cache.RedisCacheService, limit int64, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var identifier string
		if userUuid, exists := c.Get("user_uuid"); exists {
			identifier = fmt.Sprintf("user:%s", userUuid.(string))
		} else {
			identifier = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		key := fmt.Sprintf("rate_limit:%s:%s", identifier, c.FullPath())

		count, err := rdb.Incr(key, ttl)
		if err != nil {
			log.Printf("[RateLimit] Redis error: %v, skipping check", err)
			c.Next()
			return
		}

		if count > limit {
			log.Printf("[RateLimit] BLOCKING %s - Request count: %d/%d", key, count, limit)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}

		log.Printf("[RateLimit] ALLOWING %s - Request count: %d/%d", key, count, limit)
		c.Next()
	}
}
