package app

import (
	"nineshop-be/internal/handler"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type UserModule struct {
	route routes.Route
}

func NewUserModule(userService service.UserService) *UserModule {
	userHandler := handler.NewUserHandler(userService)
	userRoute := routes.NewUserRoute(userHandler)
	return &UserModule{
		route: userRoute,
	}
}
func (m *UserModule) Route() routes.Route {
	return m.route
}
