package dto

import (
	"nineshop-be/internal/db/sqlc"
)

type AddToCartParams struct {
	VariantID int `json:"variant_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

func MapParamToSqlcCart(input AddToCartParams) sqlc.AddToCartParams {
	return sqlc.AddToCartParams{
		VariantID: int32(input.VariantID),
		Quantity:  IntToInt32(input.Quantity),
	}
}

type CartRes struct {
	ID        int64  `json:"id"`
	UserID    int32  `json:"user_id"`
	VariantID int32  `json:"variant_id"`
	Quantity  *int32 `json:"quantity"`
}

type DeleteItemInCartParams struct {
	VariantID int32 `json:"variant_id" binding:"required"`
}

type UpdateAllCartParam struct {
	Items    []CartItems `json:"items" binding:"required"`
	Shipping string      `json:"shipping_method" binding:"required"`
}

type CartItems struct {
	VariantID int     `json:"variant_id" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
}

type CartResponse struct {
	Data          []sqlc.Cart `json:"data"`
	Subtotal      float64     `json:"subtotal"`
	Total         float64     `json:"total"`
	ShippingPrice float64     `json:"shipping_price"`
	Tax           float64     `json:"tax"`
}
