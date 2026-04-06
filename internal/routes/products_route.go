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
		product.GET("get-by-filter", pr.handler.GetAllByFilter)
		product.POST("/add", middleware.RoleMiddleware("admin"), pr.handler.AddProduct)
		product.GET("/list-variant/:product_id", pr.handler.GetListVariantByPid)
		product.GET("/top-3-thumbnail", pr.handler.GetTop3Thumbnail)
		product.GET("/get-by-slug/:slug", pr.handler.GetBySlug)
		product.GET("get-all", pr.handler.GetListProducts)
	}
}
