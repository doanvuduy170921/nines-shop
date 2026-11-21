package handler

import (
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

// AddToCart godoc
// @Summary      Add item to cart
// @Description  Add a product variant with specified quantity to the authenticated user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.AddToCartParams true "Cart item details"
// @Success      200 {object} dto.CartRes "Item added successfully"
// @Failure      400 {object} utils.AppError "Validation error or invalid input"
// @Failure      401 {object} utils.AppError "Unauthorized - Missing or invalid token"
// @Router       /cart/add [post]
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

// GetCartsByUserId godoc
// @Summary      Get user's cart items
// @Description  Retrieve all items in the shopping cart for the authenticated user
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.CartRes "Cart items retrieved successfully"
// @Failure      400 {object} utils.AppError "User not found"
// @Failure      401 {object} utils.AppError "Unauthorized"
// @Router       /cart [get]
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

// DeleteItem godoc
// @Summary      Delete item from cart
// @Description  Remove a specific item/variant from the user's shopping cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.DeleteItemInCartParams true "Item identification to delete"
// @Success      204 {string} string "No Content - Item deleted successfully"
// @Failure      400 {object} utils.AppError "Validation error"
// @Failure      401 {object} utils.AppError "Unauthorized"
// @Router       /cart/item [delete]
func (ch *CartHandler) DeleteItem(ctx *gin.Context) {
	var input dto.DeleteItemInCartParams
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}
	err := ch.service.DeleteItem(ctx, input)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusNoContent, nil, "Delete Item Successfully")
}

// UpdateAllCart godoc
// @Summary      Update entire cart and recalculate totals
// @Description  Update quantities or items in bulk and return cart summary including subtotal, tax, shipping, and total prices
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.UpdateAllCartParam true "Bulk update cart parameters"
// @Success      200 {object} dto.CartResponse "Cart updated successfully with calculated totals"
// @Failure      400 {object} utils.AppError "Validation error"
// @Failure      401 {object} utils.AppError "Unauthorized"
// @Router       /cart/update [put]
func (ch *CartHandler) UpdateAllCart(ctx *gin.Context) {
	var input dto.UpdateAllCartParam
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}

	cart, subtotal, total, tax, shipPrice, err := ch.service.UpdateAllCart(ctx, input)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
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
