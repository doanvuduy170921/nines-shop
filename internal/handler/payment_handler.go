package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (ph *PaymentHandler) GetAllPayment(ctx *gin.Context) {

	payments, err := ph.service.GetAllPayment(ctx)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, payments, "Get All Payment")
}
