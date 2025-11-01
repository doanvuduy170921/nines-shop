package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type ProductImagesRoute struct {
	handler *handler.ProductImagesHandler
}

func NewProductImagesRoute(handler *handler.ProductImagesHandler) *ProductImagesRoute {
	return &ProductImagesRoute{
		handler: handler,
	}
}

func (pr *ProductImagesRoute) Register(r *gin.RouterGroup) {
	productImg := r.Group("/images")
	{
		productImg.POST("/uploads", pr.handler.UploadImage)
		productImg.POST("/uploads/:product_id", pr.handler.MultipleUploadImages)

	}
	productImg.Static("/", "./uploads")
}
