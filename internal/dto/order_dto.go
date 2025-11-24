package dto

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
