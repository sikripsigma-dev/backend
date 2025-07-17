package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SupervisorHandler struct {
	supervisorService *service.SupervisorService
}

func NewSupervisorHandler(supervisorService *service.SupervisorService) *SupervisorHandler {
	return &SupervisorHandler{supervisorService}
}

func (h *SupervisorHandler) GetStudentsBySupervisor(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userID := user.Id

	students, err := h.supervisorService.GetStudentsBySupervisor(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data mahasiswa",
		})
	}

	return c.JSON(fiber.Map{
		"students": students,
	})
}