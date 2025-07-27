package config

import (
	"fmt"
	"os"
	"strings"
	"webtoon/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		Prefork:      true,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: errorHandler,
		BodyLimit:    30 * 1024 * 1024, // 30MB
	})
}
func errorHandler(c *fiber.Ctx, err error) error {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		var values []string
		for _, fieldErr := range validationErr {
			value := fmt.Sprintf("field %s is %s %s", fieldErr.Field(), fieldErr.Tag(), fieldErr.Param())
			values = append(values, value)
		}
		arguments := strings.Join(values, ", ")
		return c.Status(400).JSON(fiber.Map{
			"error": arguments,
			"path":  c.OriginalURL(),
		})
	}
	if responseStatusException, ok := err.(*response.ErrorResponse); ok {
		return c.Status(responseStatusException.Code).JSON(fiber.Map{
			"error": responseStatusException.Message,
			"path":  c.OriginalURL(),
		})
	}
	return c.Status(500).JSON(fiber.Map{
		"error": "internal server error",
		"path":  c.OriginalURL(),
	})
}
