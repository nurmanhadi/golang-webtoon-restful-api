package handler

import (
	"webtoon/internal/domain/user/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetById(c *fiber.Ctx) error
}
type handler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &handler{userService: userService}
}
func (h *handler) GetById(c *fiber.Ctx) error {
	userId := c.Params("userId", "")
	result, err := h.userService.GetById(userId)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
