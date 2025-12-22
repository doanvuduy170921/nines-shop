package dto

import "nineshop-be/internal/db/sqlc"

// Response cho CreatePendingOrder
type CreatePendingOrderResponse struct {
	PendingOrder sqlc.PendingOrder `json:"pending_order"`
	VnPayURL     string            `json:"vnpay_url,omitempty"` // Chỉ có khi payment_method là VNPAY
}
