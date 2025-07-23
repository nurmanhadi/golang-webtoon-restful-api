package middleware

import (
	"errors"
	"slices"
	"strings"
	"webtoon/pkg/response"
	"webtoon/pkg/security"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Inject struct {
	Logger *logrus.Logger
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
