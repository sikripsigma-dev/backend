package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"net/http"

	// "strconv"

	"github.com/gofiber/fiber/v2"
)

type HeadstudyHandler struct {
	headstudyService service.HeadstudyService
}

func NewHeadstudyHandler(headstudyService service.HeadstudyService) *HeadstudyHandler {
	return &HeadstudyHandler{
		headstudyService: headstudyService,
	}
}

// GetStudentsByUniversity handles GET /api/headstudy/students
func (h *HeadstudyHandler) GetStudentsByUniversity(c *fiber.Ctx) error {
	// Ambil user data dari middleware auth
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	// Pastikan UniversityID tidak 0
	if user.Headstudy.UniversityID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "university_id is missing")
	}

	students, err := h.headstudyService.GetStudentsByUniversity(user.Headstudy.UniversityID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get students")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success get students",
		"data":    students,
	})

}
