package handler

import (
	"webtoon/internal/domain/dto"
	"webtoon/internal/domain/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ChapterHandler interface {
	AddChapter(c *fiber.Ctx) error
	UpdateChapter(c *fiber.Ctx) error
	GetChapterByIdAndNumber(c *fiber.Ctx) error
	RemoveChapter(c *fiber.Ctx) error
}
type chapterHandler struct {
	chapterService service.ChapterService
}

func NewChapterHandler(chapterService service.ChapterService) ChapterHandler {
	return &chapterHandler{chapterService: chapterService}
}
func (h *chapterHandler) AddChapter(c *fiber.Ctx) error {
	request := new(dto.ChapterAddRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.chapterService.AddChapter(request); err != nil {
		return err
	}
	return response.Success(c, 201, "OK")
}
func (h *chapterHandler) UpdateChapter(c *fiber.Ctx) error {
	chapterId := c.Params("chapterId")
	request := new(dto.ChapterUpdateRequest)
	if err := c.BodyParser(&request); err != nil {
		return response.Exception(400, "error parse json")
	}
	if err := h.chapterService.UpdateChapter(chapterId, *request); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
func (h *chapterHandler) GetChapterByIdAndNumber(c *fiber.Ctx) error {
	chapterId := c.Params("chapterId")
	number := c.Query("number", "1")
	result, err := h.chapterService.GetByIdAndNumber(chapterId, number)
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
func (h *chapterHandler) RemoveChapter(c *fiber.Ctx) error {
	chapterId := c.Params("chapterId")
	if err := h.chapterService.Remove(chapterId); err != nil {
		return err
	}
	return response.Success(c, 200, "OK")
}
