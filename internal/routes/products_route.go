package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/middleware"
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
		product.POST("/create", middleware.RoleMiddleware("admin"), pr.handler.CreateProduct)
		product.GET("get-by-filter", pr.handler.GetAllByFilter)
		product.GET("/get-cate-id/:id", middleware.RoleMiddleware("admin"), pr.handler.GetAllByCateId)
		product.GET("/get-by-slug/:slug", pr.handler.GetAllBySlug)
		product.GET("/get-images-by-slug/:slug", pr.handler.GetImagesBySlug)
	}

}
