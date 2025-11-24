package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type PaymentRoute struct {
	handler *handler.PaymentHandler
}

func NewPaymentRoute(handler *handler.PaymentHandler) *PaymentRoute {
	return &PaymentRoute{
		handler: handler,
	}
}

func (cr *PaymentRoute) Register(r *gin.RouterGroup) {
	payment := r.Group("/payment")
	{
		payment.GET("/get-all", cr.handler.GetAllPayment)
	}

}
