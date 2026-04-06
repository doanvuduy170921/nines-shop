package app

import (
	"nineshop-be/internal/handler"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type AuthModule struct {
	route routes.Route
}

func NewAuthModule(authService service.AuthService, userService service.UserService) *AuthModule {
	authHandler := handler.NewAuthHandler(authService, userService)
	authRoute := routes.NewAuthRoute(authHandler)
	return &AuthModule{
		route: authRoute,
	}
}
func (m *AuthModule) Route() routes.Route {
	return m.route
}
