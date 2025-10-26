package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"nineshop-be/internal/utils"
	"time"
)

type RedisConfig struct {
	Addr     string
	Username string

	Password string
	DB       int
}

func NewRedisConfig() *redis.Client {
	cfg := RedisConfig{
		Addr:     utils.GetEnv("REDIS_HOST", "localhost:6379"),
		Username: utils.GetEnv("REDIS_USER", ""),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatal("⛔ Failed to connect to redis")
	}
	log.Println("✅ Connected to redis")
	return client
}
