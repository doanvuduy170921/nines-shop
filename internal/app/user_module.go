package app

import (
	"github.com/redis/go-redis/v9"
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"nineshop-be/pkg/email"
)

type UserModule struct {
	route routes.Route
}

func NewUserModule(redis *redis.Client, mail email.EmailService) *UserModule {
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo, redis, mail)
	userHandler := handler.NewUserHandler(userService)
	userRoute := routes.NewUserRoute(userHandler)
	return &UserModule{
		route: userRoute,
	}
}
func (m *UserModule) Route() routes.Route {
	return m.route
}
