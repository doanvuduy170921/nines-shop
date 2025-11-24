package config

import "nineshop-be/internal/utils"

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromName     string
	Environment  string
}

func NewEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     "587",
		SMTPUser:     utils.GetEnv("SMTP_USER", ""),
		SMTPPassword: utils.GetEnv("SMTP_PASSWORD", ""),
		FromName:     utils.GetEnv("FROM_NAME", "NineShop"),
		Environment:  utils.GetEnv("ENVIRONMENT", "development"),
	}
}
