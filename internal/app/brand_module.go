package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type BrandModule struct {
	route *routes.BrandRoute
}

func NewBrandModule() *BrandModule {
	BrandRepo := repository.NewBrandRepository(db.DB)
	BrandService := service.NewBrandService(BrandRepo)
	BrandHandler := handler.NewBrandHandler(BrandService)
	BrandRoute := routes.NewBrandRoute(BrandHandler)
	return &BrandModule{
		route: BrandRoute,
	}
}

func (pm *BrandModule) Route() routes.Route {
	return pm.route
}
