package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/middleware"
	"nineshop-be/pkg/auth"
	"nineshop-be/pkg/cache"
	"time"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoute(cache cache.RedisCacheService, r *gin.Engine, routes ...Route) { // ... dùng để bắt tất cả các Route
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api/v1")
	// dùng late limit cho các api không cần token
	api.Use(middleware.RateLimitMiddleware(cache, 1*time.Minute, 2))
	protected := api.Group("")

	tokenService := auth.NewJwtService(cache)
	protected.Use(
		middleware.AuthMiddleware(tokenService, cache))
	for _, route := range routes {
		switch route.(type) {
		case *AuthRoute, *ProductImagesRoute, *PaymentRoute:
			route.Register(api)
		default:
			route.Register(protected)
		}
	}
}
