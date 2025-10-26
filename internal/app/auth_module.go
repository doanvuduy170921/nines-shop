package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"nineshop-be/pkg/auth"
	"nineshop-be/pkg/cache"
)

type AuthModule struct {
	route routes.Route
}

func NewAuthModule(cache cache.RedisCacheService, userService service.UserService) *AuthModule {
	userRepo := repository.NewUserRepository(db.DB)
	tokenService := auth.NewJwtService(cache)
	authService := service.NewAuthService(userRepo, tokenService, cache)
	authHandler := handler.NewAuthHandler(authService, userService)
	authRoute := routes.NewAuthRoute(authHandler)
	return &AuthModule{
		route: authRoute,
	}
}
func (m *AuthModule) Route() routes.Route {
	return m.route
}
