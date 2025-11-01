package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type CartModule struct {
	route *routes.CartRoute
}

func NewCartModule() *CartModule {
	userRepo := repository.NewUserRepository(db.DB)
	cartRepo := repository.NewCartRepository(db.DB)
	cartService := service.NewCartService(cartRepo, userRepo)
	cartHandler := handler.NewCartHandler(cartService)
	cartRoute := routes.NewCartRoute(cartHandler)
	return &CartModule{
		route: cartRoute,
	}
}

func (pm *CartModule) Route() routes.Route {
	return pm.route
}
