package dto

import (
	"nineshop-be/internal/db/sqlc"
)

type AddToCartParams struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

func MapParamToSqlcCart(input AddToCartParams) sqlc.AddToCartParams {
	return sqlc.AddToCartParams{
		ProductID: int32(input.ProductID),
		Quantity:  IntToInt32(input.Quantity),
	}
}

type CartRes struct {
	ID        int64  `json:"id"`
	UserID    int32  `json:"user_id"`
	ProductID int32  `json:"product_id"`
	Quantity  *int32 `json:"quantity"`
}

type DeleteItemInCartParams struct {
	ProductID int32 `json:"product_id" binding:"required"`
}

type UpdateAllCartParam struct {
	Items    []CartItems `json:"items" binding:"required"`
	Shipping string      `json:"shipping_method" binding:"required"`
}

type CartItems struct {
	ProductId int     `json:"product_id" binding:"required"`
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
