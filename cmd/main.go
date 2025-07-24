package main

import (
	"context"
	"log"
	"webtoon/config"
)

func main() {
	ctx := context.Background()
	config.NewEnvirontment()
	logger := config.NewLogger()
	validation := config.NewValidator()
	db := config.NewMysql()
	s3 := config.NewMinio()
	app := config.NewFiber()

	config.Initialize(&config.Configuration{
		Ctx:        ctx,
		Logger:     logger,
		Validation: validation,
		DB:         db,
		S3:         s3,
		App:        app,
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("gofiber error: %s", err.Error())
	}
}
