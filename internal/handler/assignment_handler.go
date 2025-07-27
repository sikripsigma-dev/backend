package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"strconv"

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

// by student
func (h *AssignmentHandler) GetMyActiveAssignment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userID := user.Id

	assignment, err := h.assignmentService.GetActiveAssignment(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No active assignment found"})
	}

	return c.JSON(fiber.Map{"data": assignment})
}


func (h *AssignmentHandler) GetActiveAssignmentsByCompany(c *fiber.Ctx) error {
	// companyID := c.Params("company_id")
	user := c.Locals("user").(*models.User)
	companyID := user.Company.CompanyID

	assignments, err := h.assignmentService.GetActiveAssignmentsByCompany(companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data assignment"})
	}

	return c.JSON(fiber.Map{"data": assignments})
}

func (h *AssignmentHandler) UpdateAssignmentStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignment ID"})
	}

	if err := h.assignmentService.UpdateAssignmentStatus(uint(id), req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Assignment status updated"})
}
