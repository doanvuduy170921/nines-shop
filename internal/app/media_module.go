package app

import (
	"nineshop-be/internal/handler"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
)

type MediaModule struct {
	route *routes.MediaRoutes
}

func NewMediaModule() *MediaModule {
	mediaService := service.NewMediaService("./uploads", "http://localhost:8080/uploads/")
	mediaHandler := handler.NewMediaHandler(mediaService)
	mediaRoute := routes.NewMediaRoutes(mediaHandler)
	return &MediaModule{
		route: mediaRoute,
	}
}

func (m *MediaModule) Route() routes.Route {
	return m.route
}
