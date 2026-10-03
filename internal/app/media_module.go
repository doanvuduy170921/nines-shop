package app

import (
	"nineshop-be/internal/handler"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"os"
)

type MediaModule struct {
	route *routes.MediaRoutes
}

func NewMediaModule() *MediaModule {
	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8080/uploads/"
	}

	mediaService := service.NewMediaService("./uploads", publicURL)
	mediaHandler := handler.NewMediaHandler(mediaService)
	mediaRoute := routes.NewMediaRoutes(mediaHandler)
	return &MediaModule{
		route: mediaRoute,
	}
}

func (m *MediaModule) Route() routes.Route {
	return m.route
}
