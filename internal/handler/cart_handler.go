package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func (ch *CartHandler) AddToCart(ctx *gin.Context) {
	userUuid, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "User Not Found"})
		return
	}
	uuid := userUuid.(string)
	var input dto.AddToCartParams
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}

	cart, err := ch.service.AddToCart(ctx, uuid, input)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	res := dto.CartRes{
		ID:        cart.ID,
		UserID:    cart.UserID,
		VariantID: cart.VariantID,
		Quantity:  cart.Quantity,
	}
	utils.ResponseSuccess(ctx, http.StatusOK, res, "Add To Cart Success")

}
func (ch *CartHandler) GetCartsByUserId(ctx *gin.Context) {
	userUuid, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "User Not Found"})
		return
	}
	uuid := userUuid.(string)
	cart, err := ch.service.GetCartsByUserId(ctx, uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, cart, "Get Cart By User Success")
}

func (ch *CartHandler) DeleteItem(ctx *gin.Context) {
	var input dto.DeleteItemInCartParams
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}
	err := ch.service.DeleteItem(ctx, input)
	if err != nil {
		utils.ResponseError(ctx, err)
	}
	utils.ResponseSuccess(ctx, http.StatusNoContent, nil, "Delete Item Successfully")

}

func (ch *CartHandler) UpdateAllCart(ctx *gin.Context) {
	var input dto.UpdateAllCartParam
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}

	cart, subtotal, total, tax, shipPrice, err := ch.service.UpdateAllCart(ctx, input)
	if err != nil {
		utils.ResponseError(ctx, err)
	}
	res := dto.CartResponse{
		Data:          cart,
		Subtotal:      subtotal,
		Total:         total,
		ShippingPrice: shipPrice,
		Tax:           tax,
	}

	utils.ResponseSuccess(ctx, http.StatusOK, res, "Update All Cart Success")

}
