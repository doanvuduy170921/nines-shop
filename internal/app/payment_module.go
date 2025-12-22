package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type PaymentModule struct {
	route *routes.PaymentRoute
}

func NewPaymentModule() *PaymentModule {
	PaymentRepo := repository.NewPaymentMethodRepository(db.DB)
	PaymentService := service.NewPaymentService(PaymentRepo)
	
	PaymentHandler := handler.NewPaymentHandler(PaymentService)
	PaymentRoute := routes.NewPaymentRoute(PaymentHandler)
	return &PaymentModule{
		route: PaymentRoute,
	}
}

func (pm *PaymentModule) Route() routes.Route {
	return pm.route
}
