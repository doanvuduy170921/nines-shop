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
