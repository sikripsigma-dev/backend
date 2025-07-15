package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Handler struct untuk Weekly Report
type Handler struct {
	weeklyReportService service.WeeklyReportService
}

// NewWeeklyReportHandler membuat instance baru dari Handler
func NewWeeklyReportHandler(
	weeklyReportService service.WeeklyReportService,
) *Handler {
	return &Handler{
		weeklyReportService: weeklyReportService,
	}
}

func (h *Handler) SubmitWeeklyReport(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal membaca form data",
			"error":   err.Error(),
		})
	}

	user, ok := c.Locals("user").(*models.User)
	if !ok || user == nil || user.Id == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	input := dto.CreateWeeklyReportDTO{
		Week:      parseInt(form.Value["week"]),
		Progress:  parseString(form.Value["progress"]),
		Plans:     parseString(form.Value["plans"]),
		Mood:      parseInt(form.Value["mood"]),
		Notes:     parseString(form.Value["notes"]),
		StartDate: parseString(form.Value["start_date"]),
		EndDate:   parseString(form.Value["end_date"]),
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
	}

	files := form.File["files"]
	uploadedFiles := []string{}

	for _, file := range files {
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		savePath := fmt.Sprintf("./public/uploads/%s", filename)

		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menyimpan file",
				"error":   err.Error(),
			})
		}
		uploadedFiles = append(uploadedFiles, filename)
	}

	input.Files = uploadedFiles

	if err := h.weeklyReportService.SubmitReport(user.Id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan laporan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Laporan berhasil dikirim",
	})
}

func (h *Handler) GetWeeklyReports(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	reports, err := h.weeklyReportService.GetReports(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil laporan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil laporan",
		"data":    reports,
	})
}

func (h *Handler) GetWeeklyReportsForSupervisor(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	if user.RoleId != 4 || user.Supervisor == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses hanya untuk dosen pembimbing",
		})
	}

	universityID := user.Supervisor.UniversityID

	reports, err := h.weeklyReportService.GetReportsBySupervisor(universityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil laporan",
		})
	}

	var response []fiber.Map
	for _, r := range reports {
		initial := ""
		nameParts := strings.Split(r.Student.Name, " ")
		for _, part := range nameParts {
			if len(part) > 0 {
				initial += strings.ToUpper(part[:1])
			}
		}

		response = append(response, fiber.Map{
			"id": r.ID,
			"student": fiber.Map{
				"name":     r.Student.Name,
				"nim":      r.Student.Nim,
				"avatar":   r.Student.Image,
				"initials": initial,
			},
			"week":         r.Week,
			"date_range":   r.StartDate.Format("02-01-2006") + " to " + r.EndDate.Format("02-01-2006"),
			"progress":     r.Progress,
			"nextWeekPlan": r.Plans,
			"mood":         r.Mood,
			"status":       r.Status,
			"feedback":     r.Notes,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil laporan",
		"data":    response,
	})
}

func parseInt(values []string) int {
	if len(values) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(values[0])
	return n
}

func parseString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
