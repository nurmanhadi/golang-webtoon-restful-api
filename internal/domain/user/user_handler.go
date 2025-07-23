package user

import (
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetById(c *fiber.Ctx) error
	UpdateUsername(c *fiber.Ctx) error
	UploadAvatar(c *fiber.Ctx) error
}
type handler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
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
	request := new(UserUpdateUsernameRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.userService.UpdateUsername(userId, *request); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
func (h *handler) UploadAvatar(c *fiber.Ctx) error {
	userId := c.Params("userId", "")
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return response.Exception(400, "avatar is required")
	}
	if err := h.userService.UploadAvatar(userId, *avatar); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
