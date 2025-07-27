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
	ChapterHandler    handler.ChapterHandler
	ContentHandler    handler.ContentHandler
}

func (i *Init) Setup(app *fiber.App) {
	api := app.Group("/api")

	// public
	api.Get("/comics", i.ComicHandler.GetAll)
	api.Get("/comics/:comicId", i.ComicHandler.GetById)
	api.Get("/comics/:comicId/chapters/:chapterId", i.ChapterHandler.GetChapterByIdAndNumber)

	api.Get("/search", i.ComicHandler.Search)

	api.Get("/genres", i.GenreHandler.GetAll)
	api.Get("/genres", i.GenreHandler.GetAll)
	api.Get("/genres/:genreId", i.GenreHandler.GetById)

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

	// chapters
	chapter := comic.Group("/:comicId/chapters")
	chapter.Post("/", i.ChapterHandler.AddChapter)
	chapter.Put("/:chapterId", i.ChapterHandler.UpdateChapter)
	chapter.Delete("/:chapterId", i.ChapterHandler.RemoveChapter)

	// contents
	content := chapter.Group("/:chapterId/contents")
	content.Post("/", i.ContentHandler.AddBulkContent)
	content.Delete("/:contentId", i.ContentHandler.RemoveContent)

	// genres
	genre := api.Group("/genres", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN)}))
	genre.Post("/", i.GenreHandler.AddGenre)
	genre.Delete("/:genreId", i.GenreHandler.Remove)

	// comic genre
	comicGenre := api.Group("/comic-genre", i.Middleware.JwtValidation(), i.Middleware.RequireRole([]string{string(role.ADMIN)}))
	comicGenre.Post("/", i.ComicGenreHandler.AddComicGenre)
	comicGenre.Delete("/comics/:comicId/genres/:genreId", i.ComicGenreHandler.RemoveComicGenreByComicIdAndGenreId)
}
