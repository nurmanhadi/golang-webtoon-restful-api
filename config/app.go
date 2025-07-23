package config

import (
	"webtoon/internal/domain/auth"
	"webtoon/internal/domain/user"

	"webtoon/internal/infrastructure/rest/middleware"
	"webtoon/internal/infrastructure/rest/routes"
	"webtoon/internal/infrastructure/storage/mysql"

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
	authRepo := mysql.NewAuthRepository(conf.DB)
	userRepo := mysql.NewUserRepository(conf.DB)

	// service
	authServ := auth.NewAuthService(conf.Logger, conf.Validation, authRepo)
	userServ := user.NewUserService(conf.Logger, conf.Validation, userRepo)

	// handler
	authHand := auth.NewAuthHandler(authServ)
	userHand := user.NewUserHandler(userServ)

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
