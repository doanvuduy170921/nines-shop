package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"nineshop-be/internal/config"
)

// GmailService implement EmailService
type GmailService struct {
	config config.EmailConfig
}

// NewGmailService khởi tạo email service
func NewGmailService(config config.EmailConfig) EmailService {
	return &GmailService{config: config}
}

// SendOTPEmail gửi mã OTP xác nhận đơn hàng
func (s *GmailService) SendOTPEmail(to, otp string) error {
	subject := "Mã xác nhận đơn hàng của bạn"
	body, err := s.renderTemplate("otp", map[string]interface{}{
		"OTP":       otp,
		"ExpiresIn": "5 phút",
	})
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return s.send(to, subject, body)
}

// send - hàm private thực hiện gửi email
func (s *GmailService) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	message := s.buildMessage(to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	// Kết nối bình thường trước (plain TCP)
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Quit()

	// Nếu server hỗ trợ STARTTLS -> nâng cấp kết nối lên TLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: s.config.SMTPHost}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Xác thực
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Gửi mail
	if err := client.Mail(s.config.SMTPUser); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

// buildMessage xây dựng email message với headers
func (s *GmailService) buildMessage(to, subject, body string) string {
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.SMTPUser)

	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return message
}

// renderTemplate render HTML template cho email
func (s *GmailService) renderTemplate(templateName string, data map[string]interface{}) (string, error) {
	templates := getEmailTemplates()

	tmpl, exists := templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	t, err := template.New(templateName).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getEmailTemplates trả về các email templates
func getEmailTemplates() map[string]string {
	return map[string]string{
		"otp": `
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 10px;">
        <h2 style="color: #2c3e50; text-align: center;">Mã xác nhận đơn hàng</h2>
        <p style="font-size: 16px; color: #555;">Mã OTP của bạn là:</p>
        <div style="background-color: #e8f4f8; padding: 20px; text-align: center; border-radius: 5px; margin: 20px 0;">
            <h1 style="color: #3498db; font-size: 36px; margin: 0; letter-spacing: 5px;">{{.OTP}}</h1>
        </div>
        <p style="color: #e74c3c; font-size: 14px;">⚠️ Mã này sẽ hết hạn sau {{.ExpiresIn}}</p>
        <p style="color: #7f8c8d; font-size: 13px;">Nếu bạn không yêu cầu mã này, vui lòng bỏ qua email.</p>
    </div>
</body>
</html>`,

		"password_reset": `
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 10px;">
        <h2 style="color: #2c3e50;">Đặt lại mật khẩu</h2>
        <p style="font-size: 16px; color: #555;">Bạn đã yêu cầu đặt lại mật khẩu. Nhấn vào nút bên dưới để tiếp tục:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.ResetLink}}" style="background-color: #3498db; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">Đặt lại mật khẩu</a>
        </div>
        <p style="color: #e74c3c; font-size: 14px;">⚠️ Link này sẽ hết hạn sau {{.ExpiresIn}}</p>
        <p style="color: #7f8c8d; font-size: 13px;">Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này.</p>
    </div>
</body>
</html>`,

		"order_confirmation": `
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 10px;">
        <h2 style="color: #27ae60; text-align: center;">✓ Đơn hàng đã được xác nhận</h2>
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <p style="margin: 10px 0;"><strong>Mã đơn hàng:</strong> #{{.OrderID}}</p>
            <p style="margin: 10px 0;"><strong>Sản phẩm:</strong> {{.ProductName}}</p>
            <p style="margin: 10px 0;"><strong>Tổng tiền:</strong> {{.Amount}} VNĐ</p>
            <p style="margin: 10px 0;"><strong>Ngày đặt:</strong> {{.OrderDate}}</p>
        </div>
        <p style="color: #555;">Cảm ơn bạn đã mua hàng! Chúng tôi sẽ xử lý đơn hàng và giao đến bạn sớm nhất.</p>
    </div>
</body>
</html>`,

		"welcome": `
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 10px;">
        <h2 style="color: #2c3e50; text-align: center;">Chào mừng {{.Username}}! 🎉</h2>
        <p style="font-size: 16px; color: #555;">Cảm ơn bạn đã đăng ký tài khoản. Chúng tôi rất vui khi có bạn!</p>
        <p style="font-size: 16px; color: #555;">Bắt đầu khám phá các tính năng của chúng tôi ngay hôm nay.</p>
    </div>
</body>
</html>`,
	}
}
