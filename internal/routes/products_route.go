package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type ProductRoute struct {
	handler *handler.ProductHandler
}

func NewProductRoute(handler *handler.ProductHandler) *ProductRoute {
	return &ProductRoute{
		handler: handler,
	}
}

func (pr *ProductRoute) Register(r *gin.RouterGroup) {
	product := r.Group("/product")
	{
		product.POST("/create", pr.handler.CreateProduct)
		product.GET("get-by-filter", pr.handler.GetAllByFilter)
	}

}
