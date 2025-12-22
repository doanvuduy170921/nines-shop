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

var AllowedTransactions = map[string][]string{
	OrderStatusPending: {
		OrderStatusProcessing,
		OrderStatusCancelled,
	},
	OrderStatusProcessing: {
		OrderStatusShipping,
		OrderStatusCancelled,
	},
	OrderStatusShipping: {
		OrderStatusDelivered,
	},
	OrderStatusDelivered: {
		OrderStatusRefunded, // nếu user muốn hoàn tiền.
	},
}

// hàm kiểm tra nhảy cóc (vd: pending --> shipping )
func CanTransaction(current, next string) bool {
	allowed, ok := AllowedTransactions[current]
	if !ok {
		return false
	}
	for _, v := range allowed {
		if v == next {
			return true
		}
	}
	return false
}

var OrderStatusNotes = map[string]string{
	OrderStatusPending:    "Order created and awaiting confirmation",
	OrderStatusProcessing: "Order is being prepared",
	OrderStatusShipping:   "Order was handed to shipping unit",
	OrderStatusDelivered:  "Order delivered successfully",
	OrderStatusCancelled:  "Order cancelled",
	OrderStatusRefunded:   "Order refunded to customer",
}
