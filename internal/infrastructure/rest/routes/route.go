package routes

import (
	"webtoon/internal/domain/auth"
	"webtoon/internal/domain/comic"
	"webtoon/internal/domain/user"
	"webtoon/internal/infrastructure/rest/middleware"
	"webtoon/pkg/role"

	"github.com/gofiber/fiber/v2"
)

type Init struct {
	Middleware   *middleware.Inject
	AuthHandler  auth.AuthHandler
	UserHandler  user.UserHandler
	ComicHandler comic.ComicHandler
}

func (i *Init) Setup(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", i.AuthHandler.Register)
	auth.Post("/login", i.AuthHandler.Login)

	user := api.Group("/users", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN), string(role.USER)}))
	user.Get("/:userId", i.UserHandler.GetById)
	user.Put("/:userId", i.UserHandler.UpdateUsername)
	user.Put("/:userId/upload", i.UserHandler.UploadAvatar)

	comic := api.Group("/comics", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN), string(role.USER)}))
	comic.Post("/", i.ComicHandler.AddComic)

}
