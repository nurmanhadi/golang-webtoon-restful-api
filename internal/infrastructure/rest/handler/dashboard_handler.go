package handler

import (
	"webtoon/internal/domain/service"
	"webtoon/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler interface {
	Summary(c *fiber.Ctx) error
}
type dashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) DashboardHandler {
	return &dashboardHandler{
		dashboardService: dashboardService,
	}
}
func (h *dashboardHandler) Summary(c *fiber.Ctx) error {
	result, err := h.dashboardService.Summary()
	if err != nil {
		return err
	}
	return response.Success(c, 200, result)
}
