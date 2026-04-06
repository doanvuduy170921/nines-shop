package vnpay

const (
	Version     = "2.1.0"
	Command     = "pay"
	TmnCode     = "PJZZ1SKG"
	HashSecret  = "6KA6JHSAI6OL5OYLCTA6PCJN7ZX6JSY7"
	PaymentURL  = "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"
	Locale      = "vn"
	CurrCode    = "VND"
	TxRefPrefix = "NINESHOP" // mày để là NINESHOP123 → sẽ cắt prefix này
)

const USD_TO_VND_RATE = 25000 // Tỷ giá USD -> VND
