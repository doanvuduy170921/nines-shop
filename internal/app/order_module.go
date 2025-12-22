package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type OrderModule struct {
	route *routes.OrderRoute
}

func NewOrderModule() *OrderModule {
	OrderStatusHis := repository.NewOrderStatusHistory(db.DB)
	OrderRepo := repository.NewOrderRepository(db.DB)
	UserRepo := repository.NewUserRepository(db.DB)
	OrderService := service.NewOrderService(OrderRepo, OrderStatusHis, UserRepo)
	OrderHandler := handler.NewOrderHandler(OrderService)
	OrderRoute := routes.NewOrderRoute(OrderHandler)
	return &OrderModule{
		route: OrderRoute,
	}
}

func (pm *OrderModule) Route() routes.Route {
	return pm.route
}
