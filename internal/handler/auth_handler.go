package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userResp, err := h.authService.Register(req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": userResp,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	token, err := h.authService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Set token ke cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func (h *AuthHandler) GetUserData(c *fiber.Ctx) error {
	
	user := c.Locals("user").(*models.User)

	userResponse := fiber.Map{
		"id":    user.Id,
		"nim":   user.Nim,
		"name":  user.Name,
		"phone": user.Phone,
		"email": user.Email,
		"role":  user.RoleId,
		"image": user.Image,
	}

	if user.Company != nil {
		userResponse["company"] = fiber.Map{
			"id":       user.Company.CompanyID,
			// "name":     user.Company.CompanyName,
			"division": user.Company.Division,
		}
	}

	if user.Student != nil {
		userResponse["Student"] = fiber.Map{
			"univ_id": user.Student.UniversityID,
			"univ_name": user.Student.University.Name,
			"nim": user.Student.Nim,
			"jurusan": user.Student.Jurusan,
			"gpa": user.Student.Gpa,
		}
	}

	if user.Supervisor != nil {
		userResponse["Supervisor"] = fiber.Map{
			"univ_id": user.Supervisor.UniversityID,
			"univ_name": user.Supervisor.University.Name,
			"nidn": user.Supervisor.Nidn,
		}
	}

	return c.JSON(fiber.Map{
		"message": "User data retrieved successfully",
		"user":    userResponse,
	})
}


func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "token",
		Value: "",
		Expires: time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure: false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}