package routes

import (
	auth "webtoon/internal/domain/auth/handler"
	user "webtoon/internal/domain/user/handler"
	"webtoon/internal/infrastructure/rest/middleware"
	"webtoon/pkg/role"

	"github.com/gofiber/fiber/v2"
)

type Init struct {
	Middleware  *middleware.Inject
	AuthHandler auth.AuthHandler
	UserHandler user.UserHandler
}

func (i *Init) Setup(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", i.AuthHandler.Register)
	auth.Post("/login", i.AuthHandler.Login)

	user := api.Group("/users", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN), string(role.USER)}))
	user.Get("/:userId", i.UserHandler.GetById)
}
