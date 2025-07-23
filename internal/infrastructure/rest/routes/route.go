package routes

import (
	auth "webtoon/internal/domain/auth/handler"
	user "webtoon/internal/domain/user/handler"

	"github.com/gofiber/fiber/v2"
)

type Init struct {
	AuthHandler auth.AuthHandler
	UserHandler user.UserHandler
}

func (i *Init) Setup(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", i.AuthHandler.Register)
	auth.Post("/login", i.AuthHandler.Login)

	user := api.Group("/users")
	user.Get("/:userId", i.UserHandler.GetById)
}
