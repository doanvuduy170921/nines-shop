package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
	"strconv"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (nh *OrderHandler) UpdateStatusForUser(c *gin.Context) {
	var input dto.UpdateStatusForUserParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}
	if err := nh.service.UpdateStatusForUser(c, input); err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusNoContent, nil, "UpdateStatusForUser OK")

}

func (nh *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := nh.service.GetAllOrders(c)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, orders, "Get All Orders successfully")
}

func (nh *OrderHandler) GetOrderDetailById(c *gin.Context) {
	idStr := c.Param("id")
	id := utils.StringToInt32(idStr)

	order, err := nh.service.GetOrderDetailById(c, id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, order, "Get Order Detail successfully")
}

func (nh *OrderHandler) GetAllStatusByOrderId(c *gin.Context) {
	idStr := c.Param("id")
	id := utils.StringToInt32(idStr)
	status, err := nh.service.GetAllStatusByOrderId(c, id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, status, "Get All Status By Order ID successfully")
}

func (nh *OrderHandler) GetAllStatusByOrderIdV2(c *gin.Context) {
	idStr := c.Param("id")
	id := utils.StringToInt32(idStr)
	status, err := nh.service.GetAllStatusByOrderIdV2(c, id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, status, "Get All Status By Order ID successfully")
}

func (nh *OrderHandler) GetListOrderDetailByUserId(c *gin.Context) {
	searchStr := c.Query("search")
	status := c.Query("status")
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit := validation.ConvertAndValidateIntParam(c, limitStr, 10)
	page := validation.ConvertAndValidateIntParam(c, pageStr, 1)

	var search *int32
	if searchStr != "" {
		v, _ := strconv.Atoi(searchStr)
		tmp := int32(v)
		search = &tmp
	}

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	listOrders, err := nh.service.GetListOrdersDetailByUserId(c, search, statusPtr, limit, page)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, listOrders, "Get List Orders Detail successfully")
}

func (nh *OrderHandler) ViewOrderDetailForMyOrder(c *gin.Context) {
	orderIdStr := c.Param("order_id")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	listOrderDetails, err := nh.service.ViewDetailForMyOrder(c, int32(orderId))
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, dto.MapListViewDetailToRes(listOrderDetails), "View Detail Successfully")
}
