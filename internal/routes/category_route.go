package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
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
		category.POST("/create", cr.handler.CreateCategory)
		category.GET("/get-all", cr.handler.GetAllCategory)
	}

}
