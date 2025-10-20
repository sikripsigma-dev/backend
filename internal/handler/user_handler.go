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

	case constants.RoleSupervisor:
		var req dto.UpdateSupervisorProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		if err := h.userService.UpdateSupervisorProfile(user.Id, req); err != nil {
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

func (h *UserHandler) GetUserDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUserWithRelations(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Hindari kirim hash password
	user.Password = ""

	return c.JSON(fiber.Map{
		"user": user,
	})
}


// get all users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id") // pastikan route kamu: PUT /api/users/:id

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	if err := h.userService.UpdateUser(id, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

func (h *UserHandler) GetStudentDetail(c *fiber.Ctx) error {
	userID := c.Params("id")

	user, err := h.userService.GetStudentDetailByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Mahasiswa tidak ditemukan",
		})
	}

	// Build response
	docs := []fiber.Map{}
	for _, doc := range user.StudentDocuments {
		docs = append(docs, fiber.Map{
			"type": doc.Type,
			"path": doc.Path,
		})
	}

	return c.JSON(fiber.Map{
		"id":         user.Id,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"nim":        user.Student.Nim,
		"study_program": user.Student.Jurusan,
		"faculty":       user.Student.University.Name,
		"status":        user.Status,
		"image":         user.Image,
		"skills":        []string{}, // optionally from another table
		"portfolio":     user.Student.Description,
		"documents":     docs,
	})
}

type AssignSupervisorRequest struct {
    StudentID    string `json:"student_id" binding:"required"`
    SupervisorID string `json:"supervisor_id" binding:"required"`
}

func (h *UserHandler) AssignSupervisorToStudent(c *fiber.Ctx) error {
    var req AssignSupervisorRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

    err := h.userService.AssignSupervisorToStudent(req.StudentID, req.SupervisorID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to assign supervisor",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Supervisor assigned successfully",
    })
}
