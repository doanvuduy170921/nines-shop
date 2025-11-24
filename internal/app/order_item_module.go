package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type OrderItemModule struct {
	route *routes.OrderItemRoute
}

func NewOrderItemModule() *OrderItemModule {
	userRepo := repository.NewUserRepository(db.DB)
	OrderItemRepo := repository.NewOrderItemRepository(db.DB)
	OrderItemService := service.NewOrderItemService(OrderItemRepo, userRepo)
	OrderItemHandler := handler.NewOrderItemHandler(OrderItemService)
	OrderItemRoute := routes.NewOrderItemRoute(OrderItemHandler)
	return &OrderItemModule{
		route: OrderItemRoute,
	}
}

func (pm *OrderItemModule) Route() routes.Route {
	return pm.route
}
