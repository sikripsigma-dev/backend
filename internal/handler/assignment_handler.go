package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AssignmentHandler struct {
	assignmentService *service.AssignmentService
}

func NewAssignmentHandler(assignmentService *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{
		assignmentService: assignmentService,
	}
}

// func (h *AssignmentHandler) GetMyActiveAssignment(c *fiber.Ctx) error {
// 	user := c.Locals("user").(*models.User)

// 	userID := user.Id

// 	assignment, err := h.assignmentService.GetActiveAssignment(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Tidak ada studi kasus aktif",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"data": assignment,
// 	})
// }

func (h *AssignmentHandler) GetMyActiveAssignment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userID := user.Id

	assignment, err := h.assignmentService.GetActiveAssignment(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No active assignment found"})
	}

	return c.JSON(fiber.Map{"data": assignment})
}

