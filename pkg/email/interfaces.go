package email

// EmailService interface - dễ dàng mock khi test
type EmailService interface {
	SendOTPEmail(to, otp string) error
}
