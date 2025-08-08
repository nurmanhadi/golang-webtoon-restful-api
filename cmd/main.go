package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
	port := os.Getenv("APP_PORT")
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("gofiber error: %s", err.Error())
	}
}
