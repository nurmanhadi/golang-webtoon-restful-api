package handler

import (
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/service"
	comictype "webtoon/pkg/comic-type"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ComicHandler interface {
	AddComic(c *fiber.Ctx) error
	UpdateComic(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Remove(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	GetAllByType(c *fiber.Ctx) error
}

type comicHandler struct {
	ComicService service.ComicService
}

func NewComicHandler(ComicService service.ComicService) ComicHandler {
	return &comicHandler{ComicService: ComicService}
}
func (h *comicHandler) AddComic(c *fiber.Ctx) error {
	cover, err := c.FormFile("cover")
	if err != nil {
		return response.Exception(400, "cover is required")
	}
	title := c.FormValue("title", "")
	synopsis := c.FormValue("synopsis", "")
	author := c.FormValue("author", "")
	artist := c.FormValue("artist", "")
	comicType := c.FormValue("type", "")

	request := &dto.ComicAddRequest{
		Title:    title,
		Synopsis: synopsis,
		Author:   author,
		Artist:   artist,
		Type:     comictype.TYPE(comicType),
	}
	if err := h.ComicService.AddComic(cover, request); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
func (h *comicHandler) UpdateComic(c *fiber.Ctx) error {
	comicId := c.Params("comicId")
	cover, _ := c.FormFile("cover")
	title := c.FormValue("title", "")
	synopsis := c.FormValue("synopsis", "")
	author := c.FormValue("author", "")
	artist := c.FormValue("artist", "")
	comicType := c.FormValue("type", "")

	request := &dto.ComicUpdateRequest{
		Title:    title,
		Synopsis: synopsis,
		Author:   author,
		Artist:   artist,
		Type:     comictype.TYPE(comicType),
	}
	if err := h.ComicService.UpdateComic(comicId, cover, request); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
func (h *comicHandler) GetById(c *fiber.Ctx) error {
	comicId := c.Params("comicId")
	result, err := h.ComicService.GetById(comicId)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
func (h *comicHandler) GetAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	result, err := h.ComicService.GetAll(page, size)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
func (h *comicHandler) Remove(c *fiber.Ctx) error {
	comicId := c.Params("comicId")
	if err := h.ComicService.Remove(comicId); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
func (h *comicHandler) Search(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	result, err := h.ComicService.Search(keyword, page, size)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
func (h *comicHandler) GetAllByType(c *fiber.Ctx) error {
	comicType := c.Params("type")
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	result, err := h.ComicService.GetAllByType(comicType, page, size)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
