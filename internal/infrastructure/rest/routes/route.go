package routes

import (
	"webtoon/internal/infrastructure/rest/handler"
	"webtoon/internal/infrastructure/rest/middleware"
	"webtoon/pkg/role"

	"github.com/gofiber/fiber/v2"
)

type Init struct {
	Middleware        *middleware.Inject
	AuthHandler       handler.AuthHandler
	UserHandler       handler.UserHandler
	ComicHandler      handler.ComicHandler
	GenreHandler      handler.GenreHandler
	ComicGenreHandler handler.ComicGenreHandler
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

	// comic genre
	comicGenre := api.Group("/comic-genre", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN)}))
	comicGenre.Post("/", i.ComicGenreHandler.AddComicGenre)
	comicGenre.Delete("/comics/:comicId/genres/:genreId", i.ComicGenreHandler.RemoveComicGenreByComicIdAndGenreId)
}
