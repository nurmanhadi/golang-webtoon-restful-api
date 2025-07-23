package handler

import (
	"webtoon/internal/domain/user/dto"
	"webtoon/internal/domain/user/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetById(c *fiber.Ctx) error
	UpdateUsername(c *fiber.Ctx) error
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
func (h *handler) UpdateUsername(c *fiber.Ctx) error {
	userId := c.Params("userId", "")
	request := new(dto.UserUpdateUsernameRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.userService.UpdateUsername(userId, *request); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
