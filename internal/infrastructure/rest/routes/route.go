package routes

import (
	"webtoon/internal/domain/auth/handler"

	"github.com/gofiber/fiber/v2"
)

type Init struct {
	AuthHandler handler.AuthHandler
}

func (i *Init) Setup(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", i.AuthHandler.Register)
	auth.Post("/login", i.AuthHandler.Login)
}
