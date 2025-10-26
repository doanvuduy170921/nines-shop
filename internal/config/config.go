package config

import (
	"fmt"
	"nineshop-be/internal/utils"
)

type DatabaseConfig struct {
	DbName   string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
}

type Config struct {
	db            DatabaseConfig
	ServerAddress string
}

func NewConfig() *Config {
	return &Config{
		db: DatabaseConfig{
			DbName:   utils.GetEnv("DB_NAME", "nineshop-be"),
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5433"),
			User:     utils.GetEnv("DB_USER", "root"),
			Password: utils.GetEnv("DB_PASSWORD", ""),
			SSLMode:  utils.GetEnv("DB_SSL", "disable"),
		},
		ServerAddress: utils.GetEnv("SERVER_ADDRESS", ":8080"),
	}
}

func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.db.Host, c.db.Port, c.db.User, c.db.Password, c.db.DbName, c.db.SSLMode)
}
