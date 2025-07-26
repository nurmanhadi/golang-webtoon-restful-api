package handler

import (
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ComicGenreHandler interface {
	AddComicGenre(c *fiber.Ctx) error
	RemoveComicGenreByComicIdAndGenreId(c *fiber.Ctx) error
}
type comicGenreHandler struct {
	comicGenreService service.ComicGenreService
}

func NewComicGenreHandler(comicGenreService service.ComicGenreService) ComicGenreHandler {
	return &comicGenreHandler{comicGenreService: comicGenreService}
}
func (h *comicGenreHandler) AddComicGenre(c *fiber.Ctx) error {
	request := new(dto.ComicGenreAddRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.comicGenreService.AddComicGenre(request); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
func (h *comicGenreHandler) RemoveComicGenreByComicIdAndGenreId(c *fiber.Ctx) error {
	comicId := c.Params("comicId")
	genreId := c.Params("genreId")
	if err := h.comicGenreService.RemoveComicGenreByComicIdAndGenreId(comicId, genreId); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
