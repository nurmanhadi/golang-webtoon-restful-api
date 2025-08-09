package middleware

import (
	"errors"
	"os"
	"slices"
	"strings"
	"webtoon/pkg/response"
	"webtoon/pkg/security"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

type Inject struct {
	Logger *logrus.Logger
	App    *fiber.App
}

// setup
func (m *Inject) Setup() {

	m.App.Use(recover.New())

	m.App.Use(helmet.New())

	originUrl := os.Getenv("ORIGIN_URL")
	m.App.Use(cors.New(cors.Config{
		AllowOrigins:     originUrl,
		AllowHeaders:     "Origin, Content-Type, Authorization, Options",
		AllowMethods:     "GET, PUT, POST, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	m.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
}

// jwt verification
func (m *Inject) JwtValidation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := getTokenFromHeader(c)
		if err != nil {
			m.Logger.WithError(err).Warn("get token from header error")
			return response.Exception(401, err.Error())
		}
		claims, err := security.JwtVerify(token)
		if err != nil {
			m.Logger.WithError(err).Warn("verify token error")
			return response.Exception(401, err.Error())
		}
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
func getTokenFromHeader(c *fiber.Ctx) (string, error) {
	header := c.Get("Authorization", "")
	if header == "" {
		return "", errors.New("token null")
	}
	token := strings.Split(header, " ")
	if token[0] != "Bearer" {
		return "", errors.New("value authorization most be Bearer example 'Authorization: Bearer token'")
	}
	return token[1], nil
}

// roleValidation
func (m *Inject) RequireRole(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			m.Logger.WithField("error", role).Warn("missing or invalid role in context")
			return response.Exception(403, "you do not have permission to access this resource")
		}

		if slices.Contains(roles, role) {
			return c.Next()
		}
		m.Logger.WithField("error", role).Warn("not have permission")
		return response.Exception(403, "you do not have permission to access this resource")
	}
}
