package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/utils"
)

type CreatePendingOrderDto struct {
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Phone     string           `json:"phone"`
	Total     float64          `json:"total"`
	Subtotal  float64          `json:"subtotal"`
	Shipping  float64          `json:"shipping"`
	Tax       float64          `json:"tax"`
	PaymentId int              `json:"payment_id" binding:"required"`
	Address   string           `json:"address"  binding:"required"`
	Items     []OrderItemInput `json:"items"  binding:"required"`
}

type OrderItemInput struct {
	ProductId int     `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

type ValidateOTPParams struct {
	POrderId int    `json:"p_order_id" binding:"required"`
	Otp      string `json:"otp" binding:"required"`
}

type OTPCreateUserParam struct {
	Otp   string `json:"otp" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UpdateStatusForUserParams struct {
	Status  string `json:"status" binding:"required"`
	OrderId int    `json:"order_id" binding:"required"`
}

var shippingMethodMap = map[float64]string{
	4.99:  "Standard Delivery",
	12.99: "Express Delivery",
	0:     "Free Shipping",
}

type ViewDetailForMyOrderRes struct {
	PaymentMethodName string         `json:"payment_method_name"`
	ShippingMethod    string         `json:"shipping_method_name"`
	ShippingPrice     pgtype.Numeric `json:"shipping_price"`
	ProductName       string         `json:"product_name"`
	Quantity          int32          `json:"quantity"`
	Price             pgtype.Numeric `json:"price"`
	AmountItem        *int32         `json:"amount_item"`
	Subtotal          pgtype.Numeric `json:"subtotal"`
	Tax               pgtype.Numeric `json:"tax"`
	TotalAmount       pgtype.Numeric `json:"total_amount"`
	UserName          string         `json:"user_name"`
	Address           string         `json:"address"`
	Phone             string         `json:"phone"`
	ProductThumbnail  *string        `json:"product_thumbnail"`
}

func MapViewDetailToRes(input sqlc.ViewDetailForMyOrderRow) ViewDetailForMyOrderRes {
	shippingPrice, _ := utils.NumericToFloat64(input.ShippingPrice)
	return ViewDetailForMyOrderRes{
		PaymentMethodName: input.PaymentMethodName,
		ShippingPrice:     input.ShippingPrice,
		ShippingMethod:    shippingMethodMap[shippingPrice],
		ProductName:       input.ProductName,
		Quantity:          input.Quantity,
		Price:             input.Price,
		AmountItem:        input.AmountItem,
		Subtotal:          input.Subtotal,
		Tax:               input.Tax,
		TotalAmount:       input.TotalAmount,
		UserName:          input.UserName,
		Address:           input.Address,
		Phone:             input.Phone,
		ProductThumbnail:  input.ProductThumbnail,
	}
}

func MapListViewDetailToRes(input []sqlc.ViewDetailForMyOrderRow) []ViewDetailForMyOrderRes {
	listResult := make([]ViewDetailForMyOrderRes, len(input))
	for k, v := range input {
		listResult[k] = MapViewDetailToRes(v)
	}
	return listResult

}
