package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type ProductModule struct {
	route *routes.ProductRoute
}

func NewProductModule() *ProductModule {
	ProductRepo := repository.NewProductRepository(db.DB)
	ProductService := service.NewProductService(ProductRepo)
	ProductHandler := handler.NewProductHandler(ProductService)
	ProductRoute := routes.NewProductRoute(ProductHandler)
	return &ProductModule{
		route: ProductRoute,
	}
}

func (pm *ProductModule) Route() routes.Route {
	return pm.route
}
