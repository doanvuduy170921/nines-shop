package main

import (
	"context"
	"log"
	_ "nineshop-be/docs"
	"nineshop-be/internal/app"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db"
	"nineshop-be/internal/utils"
	"os"
	"os/signal"
	"syscall"
)

// @title Nine shop API Documentation
// @version 1.0
// @description High-performance e-commerce backend built in Golang with Gin, PostgresSQL, and Redis.
// @termsOfService  http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@nineshop.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	utils.LoadEnv()
	cfg := config.NewConfig()

	if err := db.InitDB(); err != nil {
		return err
	}
	defer db.DBPool.Close()

	app.RunMigrations(cfg.URLDNS())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application := app.NewApplication(ctx, cfg, db.DB)
	return application.Run(ctx)
}
