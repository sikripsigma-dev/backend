package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"Skripsigma-BE/internal/constants"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	switch user.RoleId {
	case constants.RoleStudent:
		var req dto.UpdateStudentProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		if err := h.userService.UpdateStudentProfile(user.Id, req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

	case constants.RoleCompany:
		var req dto.UpdateUserCompanyProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		if err := h.userService.UpdateUserCompanyProfile(user.Id, req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

	case constants.RoleAdmin:
		var req dto.UpdateAdminProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		if err := h.userService.UpdateAdminProfile(user.Id, req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized role"})
	}

	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}

// update photo profile
func (h *UserHandler) UpdateProfilePhoto(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	imageURL, err := h.userService.UpdateProfilePhoto(user.Id, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"url": imageURL,
		"message": "Profile photo updated successfully",
	})
}
