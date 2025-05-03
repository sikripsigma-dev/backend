package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)


type ApplicationHandler struct {
	applicationService *service.ApplicationService
}

func NewApplicationHandler(applicationService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService}
}

func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	var req dto.CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user := c.Locals("user").(*models.User)
	req.UserID = user.Id //get user ID from token

	application, err := h.applicationService.CreateApplication(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Applied successfully",
		"application": application,
	})
}