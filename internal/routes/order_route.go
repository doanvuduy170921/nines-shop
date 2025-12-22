package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/middleware"
)

type OrderRoute struct {
	handler *handler.OrderHandler
}

func NewOrderRoute(handler *handler.OrderHandler) *OrderRoute {
	return &OrderRoute{
		handler: handler,
	}
}

func (cr *OrderRoute) Register(r *gin.RouterGroup) {
	order := r.Group("/order")
	{
		order.PUT("/update", cr.handler.UpdateStatusForUser)
		order.GET("/get-all", middleware.RoleMiddleware("admin"), cr.handler.GetAllOrders)
		order.GET("/:id", cr.handler.GetOrderDetailById)
		order.GET("/get-tracking/:id", middleware.RoleMiddleware("admin"), cr.handler.GetAllStatusByOrderId)
		order.GET("/get-tracking-v2/:id", middleware.RoleMiddleware("admin", "customer"), cr.handler.GetAllStatusByOrderIdV2)
		order.GET("/get-all-by-user", middleware.RoleMiddleware("customer"), cr.handler.GetListOrderDetailByUserId)
		order.GET("/view-detail-for-my-order/:order_id", cr.handler.ViewOrderDetailForMyOrder)
	}
}
