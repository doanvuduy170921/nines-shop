package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
)

type PendingOrderHandler struct {
	service service.PendingOrderService
}

func NewPendingOrderHandler(service service.PendingOrderService) *PendingOrderHandler {
	return &PendingOrderHandler{
		service: service,
	}
}
func (ph *PendingOrderHandler) CreatePendingOrder(ctx *gin.Context) {
	var input dto.CreatePendingOrderDto
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}

	pOrder, err := ph.service.Create(ctx, input)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, 200, pOrder, "Create PendingOrder successfully")

}

func (ph *PendingOrderHandler) ValidateOTP(ctx *gin.Context) {
	var input dto.ValidateOTPParams
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}
	order, err := ph.service.ValidateOTP(ctx, input)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, 200, order, "Validate OTP and Create Order successfully")
}
