package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"nineshop-be/pkg/cache"
	"time"
)

func RateLimitMiddleware(rdb cache.RedisCacheService, ttl time.Duration, limit int64) gin.HandlerFunc {

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
			c.Next()
			return
		}
		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				gin.H{
					"message": "Too Many Requests,Please try again after one minute."})
			return
		}
		return
	}
}
