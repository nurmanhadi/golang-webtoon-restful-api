package routes

import (
	"webtoon/internal/domain/auth"
	"webtoon/internal/domain/comic"
	"webtoon/internal/domain/genre"
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
	GenreHandler genre.GenreHandler
}

func (i *Init) Setup(app *fiber.App) {
	api := app.Group("/api")

	// public
	api.Get("/comics", i.ComicHandler.GetAll)
	api.Get("/comics/:comicId", i.ComicHandler.GetById)

	api.Get("/search", i.ComicHandler.Search)

	// auth
	auth := api.Group("/auth")
	auth.Post("/register", i.AuthHandler.Register)
	auth.Post("/login", i.AuthHandler.Login)

	// users
	user := api.Group("/users", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN), string(role.USER)}))
	user.Get("/:userId", i.UserHandler.GetById)
	user.Put("/:userId", i.UserHandler.UpdateUsername)
	user.Put("/:userId/upload", i.UserHandler.UploadAvatar)

	//  comics
	comic := api.Group("/comics", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN)}))
	comic.Post("/", i.ComicHandler.AddComic)
	comic.Put("/:comicId", i.ComicHandler.UpdateComic)
	comic.Delete("/:comicId", i.ComicHandler.Remove)

	// genres
	genre := api.Group("/genres", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN)}))
	genre.Post("/", i.GenreHandler.AddGenre)
	genre.Get("/", i.GenreHandler.GetAll)
	genre.Delete("/:genreId", i.GenreHandler.Remove)
}
