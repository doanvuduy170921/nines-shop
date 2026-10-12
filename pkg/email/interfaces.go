package email

type EmailService interface {
	SendOTPEmail(to, otp string) error
}
