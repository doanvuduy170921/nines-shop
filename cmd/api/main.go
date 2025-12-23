package main

import (
	"log"
	"nineshop-be/internal/app"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db"
	"nineshop-be/internal/utils"
)

func main() {

	utils.LoadEnv()

	cfg := config.NewConfig()

	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
