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
		ProductID: cart.ProductID,
		Quantity:  cart.Quantity,
	}
	utils.ResponseSuccess(ctx, http.StatusOK, res, "Add To Cart Success")

}
