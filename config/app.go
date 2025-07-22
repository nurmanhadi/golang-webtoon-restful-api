package config

import (
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

	route := &routes.Init{}
	route.Setup(conf.App)
}
