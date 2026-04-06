package config

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"log"
	"nineshop-be/internal/utils"
	"time"
)

func NewRedisConfig() *redis.Client {
	addr := utils.GetEnv("REDIS_HOST", "localhost:6379")
	password := utils.GetEnv("REDIS_PASSWORD", "")
	username := utils.GetEnv("REDIS_USER", "default")
	env := utils.GetEnv("ENVIRONMENT", "development")

	opt := &redis.Options{
		Addr:         addr,
		Username:     username,
		Password:     password,
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	// Upstash production cần TLS
	if env == "production" {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	client := redis.NewClient(opt)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatal("⛔ Failed to connect to redis: ", err)
	}
	log.Println("✅ Connected to redis")
	return client
}
