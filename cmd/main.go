package main

import (
	"log"
	"webtoon/config"
)

func main() {
	config.NewEnvirontment()
	logger := config.NewLogger()
	validation := config.NewValidator()
	db := config.NewMysql()
	config.NewMinio()
	app := config.NewFiber()

	config.Initialize(&config.Configuration{
		Logger:     logger,
		Validation: validation,
		DB:         db,
		App:        app,
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("gofiber error: %s", err.Error())
	}
}
