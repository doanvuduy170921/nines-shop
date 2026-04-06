package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type MediaRoutes struct {
	handler *handler.MediaHandler
}

func NewMediaRoutes(handler *handler.MediaHandler) *MediaRoutes {
	return &MediaRoutes{
		handler: handler,
	}
}

func (mr *MediaRoutes) Register(r *gin.RouterGroup) {
	media := r.Group("/media")
	{
		media.POST("/upload", mr.handler.UploadImages)
	}

}
