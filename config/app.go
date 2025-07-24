package config

import (
	"context"
	"webtoon/internal/domain/auth"
	"webtoon/internal/domain/comic"
	"webtoon/internal/domain/user"

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

	// service
	authServ := auth.NewAuthService(conf.Logger, conf.Validation, authStore)
	userServ := user.NewUserService(conf.Logger, conf.Validation, userStore, s3Store)
	comicServ := comic.NewComicService(conf.Logger, conf.Validation, comicStore, s3Store)

	// handler
	authHand := auth.NewAuthHandler(authServ)
	userHand := user.NewUserHandler(userServ)
	comicHand := comic.NewComicHandler(comicServ)

	// middleware
	middleware := &middleware.Inject{
		Logger: conf.Logger,
	}

	route := &routes.Init{
		Middleware:   middleware,
		AuthHandler:  authHand,
		UserHandler:  userHand,
		ComicHandler: comicHand,
	}
	route.Setup(conf.App)
}
