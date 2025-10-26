package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/utils"
	"time"
)

var DB *sqlc.Queries

func InitDB() error {
	utils.LoadEnv()
	connStr := config.NewConfig().DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}
	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbPool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	DB = sqlc.New(dbPool)
	if err := dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Printf("✅ Database connection established successfully")
	return nil

}
