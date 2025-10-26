package app

import (
	"nineshop-be/internal/db"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type CategoryModule struct {
	route *routes.CategoryRoute
}

func NewCategoryModule() *CategoryModule {
	CategoryRepo := repository.NewCategoryRepository(db.DB)
	CategoryService := service.NewCategoryService(CategoryRepo)
	CategoryHandler := handler.NewCategoryHandler(CategoryService)
	CategoryRoute := routes.NewCategoryRoute(CategoryHandler)
	return &CategoryModule{
		route: CategoryRoute,
	}
}

func (pm *CategoryModule) Route() routes.Route {
	return pm.route
}
