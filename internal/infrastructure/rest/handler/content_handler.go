package handler

import (
	"webtoon/internal/domain/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ContentHandler interface {
	AddBulkContent(c *fiber.Ctx) error
}
type contentHandler struct {
	contentService service.ContentService
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{contentService: contentService}
}
func (h *contentHandler) AddBulkContent(c *fiber.Ctx) error {
	c.Params("comicId")
	chapterId := c.Params("chapterId")
	form, err := c.MultipartForm()
	if err != nil {
		response.Exception(400, "failed parse form")
	}
	contents := form.File["contents"]
	if len(contents) == 0 {
		return response.Exception(400, "no contents upload")
	}
	if len(contents) > 10 {
		return response.Exception(400, "contents max upload quantity 10")
	}
	if err := h.contentService.AddBulkContent(chapterId, contents); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
