package handler

import (
	"webtoon/internal/domain/auth/model"
	"webtoon/internal/domain/auth/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}
type handler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &handler{authService: authService}
}

func (h *handler) Register(c *fiber.Ctx) error {
	request := new(model.AuthRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.authService.Register(*request); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
func (h *handler) Login(c *fiber.Ctx) error {
	request := new(model.AuthRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	result, err := h.authService.Login(*request)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
