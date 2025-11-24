package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
)

type OrderItemHandler struct {
	service service.OrderItemService
}

func NewOrderItemHandler(service service.OrderItemService) *OrderItemHandler {
	return &OrderItemHandler{
		service: service,
	}
}

func (nh *OrderItemHandler) GetListOrderItemByUserId(c *gin.Context) {
	orderItems, err := nh.service.GetListOrderItemByUserId(c)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, orderItems, "Get List OrderItem By UserId successfully")
}
