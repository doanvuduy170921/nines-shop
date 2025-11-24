package app

import (
	"nineshop-be/internal/config"
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"nineshop-be/pkg/email"
)

type PendingOrderModule struct {
	route *routes.PendingOrderRoute
}

func NewPendingOrderModule() *PendingOrderModule {
	userRepo := repository.NewUserRepository(db.DB)
	PendingOrderRepo := repository.NewPendingOrderRepository(db.DB)
	PendingOrderItemRepo := repository.NewPendingOrderItemRepository(db.DB)
	OrderRepo := repository.NewOrderRepository(db.DB)
	OrderItemRepo := repository.NewOrderItemRepository(db.DB)
	ProductRepo := repository.NewProductRepository(db.DB)
	OrderStatusHisRepo := repository.NewOrderStatusHistory(db.DB)
	EmailConfig := config.NewEmailConfig()
	EmailService := email.NewGmailService(*EmailConfig)
	PendingOrderService := service.NewPendingOrderService(PendingOrderRepo, userRepo, PendingOrderItemRepo, EmailService, OrderRepo, OrderItemRepo, ProductRepo, OrderStatusHisRepo)
	PendingOrderHandler := handler.NewPendingOrderHandler(PendingOrderService)
	PendingOrderRoute := routes.NewPendingOrderRoute(PendingOrderHandler)
	return &PendingOrderModule{
		route: PendingOrderRoute,
	}
}

func (pm *PendingOrderModule) Route() routes.Route {
	return pm.route
}
