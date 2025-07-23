package config

import (
	"webtoon/internal/domain/auth/handler"
	"webtoon/internal/domain/auth/repository"
	"webtoon/internal/domain/auth/service"
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
	authRepo := repository.NewAuthRepository(conf.DB)

	// service
	authServ := service.NewAuthService(conf.Logger, conf.Validation, authRepo)

	// handler
	authHand := handler.NewAuthHandler(authServ)

	route := &routes.Init{
		AuthHandler: authHand,
	}
	route.Setup(conf.App)
}
