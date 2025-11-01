package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type ProductImagesModule struct {
	route *routes.ProductImagesRoute
}

func NewProductImagesModule() *ProductImagesModule {
	productRepo := repository.NewProductRepository(db.DB)
	ProductImagesRepo := repository.NewProductImagesRepository(db.DB)
	ProductImagesService := service.NewProductImagesService(ProductImagesRepo, productRepo)
	ProductImagesHandler := handler.NewProductImagesHandler(ProductImagesService)
	ProductImagesRoute := routes.NewProductImagesRoute(ProductImagesHandler)
	return &ProductImagesModule{
		route: ProductImagesRoute,
	}
}

func (pm *ProductImagesModule) Route() routes.Route {
	return pm.route
}
