package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type BrandRoute struct {
	handler *handler.BrandHandler
}

func NewBrandRoute(handler *handler.BrandHandler) *BrandRoute {
	return &BrandRoute{
		handler: handler,
	}
}

func (cr *BrandRoute) Register(r *gin.RouterGroup) {
	Brand := r.Group("/brand")
	{
		Brand.GET("/get-all", cr.handler.GetAllBrand)
	}

}
