package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type OrderItemRoute struct {
	handler *handler.OrderItemHandler
}

func NewOrderItemRoute(handler *handler.OrderItemHandler) *OrderItemRoute {
	return &OrderItemRoute{
		handler: handler,
	}
}

func (cr *OrderItemRoute) Register(r *gin.RouterGroup) {
	orderItem := r.Group("/order-item")
	{
		orderItem.GET("/get-all", cr.handler.GetListOrderItemByUserId)
	}

}
