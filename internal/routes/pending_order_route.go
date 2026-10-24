package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/middleware"
)

type PendingOrderRoute struct {
	handler *handler.PendingOrderHandler
}

func NewPendingOrderRoute(handler *handler.PendingOrderHandler) *PendingOrderRoute {
	return &PendingOrderRoute{
		handler: handler,
	}
}

func (cr *PendingOrderRoute) Register(r *gin.RouterGroup) {
	pOrder := r.Group("/pending-order")
	{
		pOrder.POST("/create", middleware.RoleMiddleware("customer"), cr.handler.CreatePendingOrder)
		pOrder.POST("/validate-otp", cr.handler.ValidateOTP)
	}

}
