package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CreatePendingOrderResponse struct {
	ID              int32              `json:"id"`
	UserID          int32              `json:"user_id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone"`
	PaymentMethodID int32              `json:"payment_method_id"`
	Address         string             `json:"address"`
	Subtotal        pgtype.Numeric     `json:"subtotal"`
	TotalAmount     pgtype.Numeric     `json:"total_amount"`
	ShippingPrice   pgtype.Numeric     `json:"shipping_price"`
	Tax             pgtype.Numeric     `json:"tax"`
	Status          string             `json:"status"`
	CreatedAt       pgtype.Timestamp   `json:"created_at"`
	ExpiresAt       pgtype.Timestamptz `json:"expires_at"`
	AmountItem      *int32             `json:"amount_item"`
	VnPayURL        string             `json:"vnpay_url,omitempty"`
}
