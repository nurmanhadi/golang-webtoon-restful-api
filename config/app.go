package config

import (
	"context"
	"webtoon/internal/domain/service"

	"webtoon/internal/infrastructure/rest/handler"
	"webtoon/internal/infrastructure/rest/middleware"
	"webtoon/internal/infrastructure/rest/routes"
	"webtoon/internal/infrastructure/storage/mysql"
	"webtoon/internal/infrastructure/storage/s3"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Configuration struct {
	Ctx        context.Context
	Logger     *logrus.Logger
	Validation *validator.Validate
	DB         *gorm.DB
	S3         *minio.Client
	App        *fiber.App
}

func Initialize(conf *Configuration) {
	// storage
	s3Store := s3.NewS3Storage(conf.Ctx, conf.S3)
	authStore := mysql.NewAuthStorage(conf.DB)
	userStore := mysql.NewUserStorage(conf.DB)
	comicStore := mysql.NewComicStorage(conf.DB)
	genreStore := mysql.NewGenreStorage(conf.DB)
	comicGenreStore := mysql.NewComicGenreStorage(conf.DB)
	chapterStore := mysql.NewChapterStorage(conf.DB)
	contentStore := mysql.NewContentStorage(conf.DB)

	// service
	authServ := service.NewAuthService(conf.Logger, conf.Validation, authStore)
	userServ := service.NewUserService(conf.Logger, conf.Validation, userStore, s3Store)
	comicServ := service.NewComicService(conf.Logger, conf.Validation, comicStore, s3Store)
	genreServ := service.NewGenreService(conf.Logger, conf.Validation, genreStore, comicGenreStore)
	comicGenreServ := service.NewComicGenreService(conf.Logger, conf.Validation, comicGenreStore, comicStore, genreStore)
	chapterServ := service.NewChapterService(conf.Logger, conf.Validation, chapterStore, comicStore)
	contentServ := service.NewContentService(conf.Logger, conf.Validation, contentStore, chapterStore, s3Store)

	// handler
	authHand := handler.NewAuthHandler(authServ)
	userHand := handler.NewUserHandler(userServ)
	comicHand := handler.NewComicHandler(comicServ)
	genreHand := handler.NewGenreHandler(genreServ)
	comicGenreHand := handler.NewComicGenreHandler(comicGenreServ)
	chapterHand := handler.NewChapterHandler(chapterServ)
	contentHand := handler.NewContentHandler(contentServ)

	// middleware
	middleware := &middleware.Inject{
		Logger: conf.Logger,
	}

	route := &routes.Init{
		Middleware:        middleware,
		AuthHandler:       authHand,
		UserHandler:       userHand,
		ComicHandler:      comicHand,
		GenreHandler:      genreHand,
		ComicGenreHandler: comicGenreHand,
		ChapterHandler:    chapterHand,
		ContentHandler:    contentHand,
	}
	route.Setup(conf.App)
}
