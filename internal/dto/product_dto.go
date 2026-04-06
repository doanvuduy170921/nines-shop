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

type ProductResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Sku              string  `json:"sku"`
	BrandName        string  `json:"brand_name"`
	CategoryName     string  `json:"category_name"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Price            float64 `json:"price"`
	DiscountPrice    float64 `json:"discount_price"`
	StockQuantity    int32   `json:"stock_quantity"`
	Status           string  `json:"status"`
	Thumbnail        string  `json:"thumbnail"`
}

type CreateCategoryParamDto struct {
	Name string `json:"name" binding:"required"`
}

func IntToInt32(i int) *int32 {
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

type AddProductRequestDto struct {
	Name             string   `json:"name" binding:"required"`
	BrandID          int      `json:"brand_id" binding:"required"`
	CategoryID       int      `json:"category_id" binding:"required"`
	Description      string   `json:"description" binding:"required"`
	ShortDescription string   `json:"short_description" binding:"required"`
	Status           string   `json:"status" binding:"required"`
	Thumbnail        string   `json:"thumbnail" binding:"required"`
	HasVariant       bool     `json:"has_variant" binding:"required"`
	Images           []string `json:"images"`
	Specifications   []struct {
		SpecKey      string `json:"spec_key" binding:"required"`
		SpecValue    string `json:"spec_value" binding:"required"`
		DisplayOrder int    `json:"display_order" binding:"required"`
	} `json:"specifications"`

	Variants []struct {
		Attributes    map[string]interface{} `json:"attributes" binding:"required"`
		Price         float64                `json:"price" binding:"required"`
		StockQuantity int32                  `json:"stock_quantity" binding:"required"`
		Image         string                 `json:"image" binding:"required"`
		IsActive      bool                   `json:"is_active" binding:"required"`
	}
}

type GetTop3TrendingRes struct {
	Images    []sqlc.GetTop3TrendingRow `json:"images"`
	Laptops   []sqlc.GetTop3TrendingRow `json:"laptops"`
	Keyboards []sqlc.GetTop3TrendingRow `json:"keyboards"`
	Screens   []sqlc.GetTop3TrendingRow `json:"screens"`
	Mouses    []sqlc.GetTop3TrendingRow `json:"mouses"`
}

type GetProductsRequest struct {
	CateName   string  `form:"cate_name"`
	Page       int32   `form:"page" binding:"required,min=1"`
	Limit      int32   `form:"limit" binding:"required,min=1,max=100"`
	SearchName string  `form:"search"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	SortBy     bool    `form:"sort_by"` // Ví dụ: price, created_at
}
