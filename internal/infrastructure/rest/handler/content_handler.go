package handler

import (
	"webtoon/internal/domain/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ContentHandler interface {
	AddBulkContent(c *fiber.Ctx) error
	RemoveContent(c *fiber.Ctx) error
}
type contentHandler struct {
	contentService service.ContentService
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{contentService: contentService}
}
func (h *contentHandler) AddBulkContent(c *fiber.Ctx) error {
	chapterId := c.Params("chapterId")
	form, err := c.MultipartForm()
	if err != nil {
		response.Exception(400, "failed parse form")
	}
	contents := form.File["contents"]
	if len(contents) == 0 {
		return response.Exception(400, "no contents upload")
	}
	if err := h.contentService.AddBulkContent(chapterId, contents); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
func (h *contentHandler) RemoveContent(c *fiber.Ctx) error {
	contentId := c.Params("contentId")
	if err := h.contentService.Remove(contentId); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
