package config

import (
	authH "webtoon/internal/domain/auth/handler"
	authR "webtoon/internal/domain/auth/repository"
	authS "webtoon/internal/domain/auth/service"
	userH "webtoon/internal/domain/user/handler"
	userR "webtoon/internal/domain/user/repository"
	userS "webtoon/internal/domain/user/service"

	"webtoon/internal/infrastructure/rest/middleware"
	"webtoon/internal/infrastructure/rest/routes"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Configuration struct {
	Logger     *logrus.Logger
	Validation *validator.Validate
	DB         *gorm.DB
	App        *fiber.App
}

func Initialize(conf *Configuration) {
	// repository
	authRepo := authR.NewAuthRepository(conf.DB)
	userRepo := userR.NewUserRepository(conf.DB)

	// service
	authServ := authS.NewAuthService(conf.Logger, conf.Validation, authRepo)
	userServ := userS.NewUserService(conf.Logger, conf.Validation, userRepo)

	// handler
	authHand := authH.NewAuthHandler(authServ)
	userHand := userH.NewUserHandler(userServ)

	// middleware
	middleware := &middleware.Inject{
		Logger: conf.Logger,
	}

	route := &routes.Init{
		Middleware:  middleware,
		AuthHandler: authHand,
		UserHandler: userHand,
	}
	route.Setup(conf.App)
}
