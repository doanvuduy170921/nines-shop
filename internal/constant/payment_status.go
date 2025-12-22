package constant

const (
	PaymentStatusUnpaid   = "UNPAID"   // COD chưa thanh toán
	PaymentStatusPaid     = "PAID"     // Đã thanh toán (VNPay/COD đã thu tiền)
	PaymentStatusRefunded = "REFUNDED" // Đã hoàn tiền (optional)
)
