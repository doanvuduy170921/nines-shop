package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/middleware"
)

type CartRoute struct {
	handler *handler.CartHandler
}

func NewCartRoute(handler *handler.CartHandler) *CartRoute {
	return &CartRoute{
		handler: handler,
	}
}

func (cr *CartRoute) Register(r *gin.RouterGroup) {
	cart := r.Group("/cart")
	{
		cart.POST("/add-to-cart", middleware.RoleMiddleware("customer"), cr.handler.AddToCart)
		cart.GET("/get-all-in-cart", middleware.RoleMiddleware("customer"), cr.handler.GetCartsByUserId)
		cart.DELETE("/delete", middleware.RoleMiddleware("customer"), cr.handler.DeleteItem)
		cart.PUT("/update-all", middleware.RoleMiddleware("customer"), cr.handler.UpdateAllCart)
	}

}
