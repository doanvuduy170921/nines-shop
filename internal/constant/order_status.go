package constant

const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipping   = "shipping"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
	OrderStatusRefunded   = "refunded"
)

var ValidOrderStatus = map[string]bool{
	OrderStatusPending:    true,
	OrderStatusProcessing: true,
	OrderStatusShipping:   true,
	OrderStatusDelivered:  true,
	OrderStatusCancelled:  true,
	OrderStatusRefunded:   true,
}

func IsValidOrderStatus(orderStatus string) bool {
	return ValidOrderStatus[orderStatus]
}
