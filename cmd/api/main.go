package main

import (
	"context"
	"log"
	"nineshop-be/internal/app"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db"
	"nineshop-be/internal/utils"
	"os"
	"os/signal"
	"syscall"
)

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
