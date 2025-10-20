package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
	// "Skripsigma-BE/internal/dto"
	"strconv"
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


// func (h *SupervisorHandler) GetAllByHeadstudy(c *fiber.Ctx) error {
// 	user := c.Locals("user").(*models.User)

// 	// Cari study program dari user kaprodi
// 	studyProgramUser  := user.Headstudy.UserID
// 	if err := h.supervisorService.DB.Where("user_id = ?", user.Id).First(&studyProgramUser).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Data kaprodi tidak ditemukan di study program user",
// 		})
// 	}

// 	// Cari semua dosen pembimbing yang berada di universitas yang sama
// 	var supervisors []models.SupervisorUser
// 	if err := h.supervisorService.DB.
// 		Preload("User").
// 		Where("university_id = ?", user.Headstudy.UniversityID).
// 		Find(&supervisors).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Gagal mengambil data dosen pembimbing",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"data": supervisors,
// 	})
// }

func (h *SupervisorHandler) GetAllByHeadstudy(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Konversi string ke uint
	userID, err := strconv.ParseUint(user.Id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID pengguna tidak valid",
		})
	}

	// Panggil service
	supervisors, err := h.supervisorService.GetAllSupervisorsByHeadstudy(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data dosen pembimbing",
		})
	}

	return c.JSON(fiber.Map{
		"data": supervisors,
	})
}

