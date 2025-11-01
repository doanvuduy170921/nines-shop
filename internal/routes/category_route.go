package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/middleware"
)

type CategoryRoute struct {
	handler *handler.CategoryHandler
}

func NewCategoryRoute(handler *handler.CategoryHandler) *CategoryRoute {
	return &CategoryRoute{
		handler: handler,
	}
}

func (cr *CategoryRoute) Register(r *gin.RouterGroup) {
	category := r.Group("/category")
	{
		category.POST("/create", middleware.RoleMiddleware("admin"), cr.handler.CreateCategory)
		category.GET("/get-all", cr.handler.GetAllCategory)
	}

}
