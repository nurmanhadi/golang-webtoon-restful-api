package comic

import (
	comictype "webtoon/pkg/comic-type"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ComicHandler interface {
	AddComic(c *fiber.Ctx) error
	UpdateComic(c *fiber.Ctx) error
}

type comicHandler struct {
	ComicService ComicService
}

func NewComicHandler(ComicService ComicService) ComicHandler {
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

	request := &ComicAddRequest{
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

	request := &ComicUpdateRequest{
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
