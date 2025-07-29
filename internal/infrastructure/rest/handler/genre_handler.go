package handler

import (
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type GenreHandler interface {
	AddGenre(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Remove(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
}

type genreHandler struct {
	genreService service.GenreService
}

func NewGenreHandler(genreService service.GenreService) GenreHandler {
	return &genreHandler{genreService: genreService}
}
func (h *genreHandler) AddGenre(c *fiber.Ctx) error {
	request := new(dto.GenreAddRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.genreService.AddGenre(request); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
func (h *genreHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.genreService.GetAll()
	if err != nil {
		return err
	}
	return response.Success(c, 201, result)
}
func (h *genreHandler) GetById(c *fiber.Ctx) error {
	genreId := c.Params("genreId")
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	result, err := h.genreService.GetById(genreId, page, size)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
func (h *genreHandler) Remove(c *fiber.Ctx) error {
	genreId := c.Params("genreId")
	if err := h.genreService.Remove(genreId); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
