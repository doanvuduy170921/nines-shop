package dto

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"nineshop-be/internal/db/sqlc"
)

type CreateProductParamDto struct {
	Name             string  `json:"name" binding:"required"`
	BrandID          int     `json:"brand_id" binding:"required"`
	CategoryID       int     `json:"category_id" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	ShortDescription string  `json:"short_description" binding:"required"`
	Price            float64 `json:"price" binding:"required"`
	DiscountPrice    float64 `json:"discount_price" binding:"required"`
	StockQuantity    int     `json:"stock_quantity" binding:"required"`
}

func MapProductDtoToParams(input CreateProductParamDto) sqlc.CreateProductParams {
	return sqlc.CreateProductParams{
		Name:             input.Name,
		BrandID:          intToInt32(input.BrandID),
		CategoryID:       intToInt32(input.CategoryID),
		Description:      &input.Description,
		ShortDescription: &input.ShortDescription,
		Price:            Float64ToPgTypeNumeric(input.Price),
		DiscountPrice:    Float64ToPgTypeNumeric(input.DiscountPrice),
		StockQuantity:    intToInt32(input.StockQuantity),
	}
}

type ProductResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Sku              string  `json:"sku"`
	BrandID          int     `json:"brand_id"`
	CategoryID       int     `json:"category_id"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Price            float64 `json:"price"`
	DiscountPrice    float64 `json:"discount_price"`
	StockQuantity    int32   `json:"stock_quantity"`
	Status           string  `json:"status"`
}

func MapSqlcProductToResponse(input sqlc.Product) ProductResponse {
	return ProductResponse{
		ID:               int(input.ID),
		Name:             input.Name,
		Slug:             input.Slug,
		Sku:              input.Sku,
		BrandID:          Int32ToInt(*input.BrandID),
		CategoryID:       Int32ToInt(*input.CategoryID),
		Description:      *input.Description,
		ShortDescription: *input.ShortDescription,
		Price:            pgTypeNumericToFloat64(input.Price),
		DiscountPrice:    pgTypeNumericToFloat64(input.DiscountPrice),
		StockQuantity:    *input.StockQuantity,
		Status:           *input.Status,
	}
}

func MapSqlcProductsToResponse(products []sqlc.Product) []ProductResponse {
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = MapSqlcProductToResponse(product) // Gán trực tiếp
	}
	return productResponses
}

type CreateCategoryParamDto struct {
	Name string `json:"name" binding:"required"`
}

func intToInt32(i int) *int32 {
	v := int32(i)
	return &v
}

func Int32ToInt(i int32) int {
	return int(i)
}

func pgTypeNumericToFloat64(input pgtype.Numeric) float64 {
	val, _ := input.Float64Value()
	if val.Valid {
		return val.Float64
	}
	return 0
}

func Float64ToPgTypeNumeric(input float64) pgtype.Numeric {
	var num pgtype.Numeric
	// chuyển float64 sang string để Scan hợp lệ
	strValue := fmt.Sprintf("%f", input)

	err := num.Scan(strValue)
	if err != nil {
		log.Println("could not convert float64 to pgtype.Numeric:", err)
		panic(err)
	}
	return num
}
