package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
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
		pOrder.POST("/create", cr.handler.CreatePendingOrder)
		pOrder.POST("/validate-otp", cr.handler.ValidateOTP)
	}

}
