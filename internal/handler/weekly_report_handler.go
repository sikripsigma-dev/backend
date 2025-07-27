package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"Skripsigma-BE/internal/util"

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
	notificationService service.NotificationService
	researchCaseService *service.ResearchCaseService
}

// NewWeeklyReportHandler membuat instance baru dari Handler
func NewWeeklyReportHandler(
	weeklyReportService service.WeeklyReportService,
	notificationService service.NotificationService,
	researchCaseService *service.ResearchCaseService,
) *Handler {
	return &Handler{
		weeklyReportService: weeklyReportService,
		notificationService: notificationService,
		researchCaseService: researchCaseService,
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
		ResearchCaseID: parseString(form.Value["research_case_id"]),
		Week:      parseInt(form.Value["week"]),
		Progress:  parseString(form.Value["progress"]),
		Plans:     parseString(form.Value["plans"]),
		Mood:      parseInt(form.Value["mood"]),
		Notes:     parseString(form.Value["notes"]),
		StartDate: parseString(form.Value["start_date"]),
		EndDate:   parseString(form.Value["end_date"]),
	}

	if input.ResearchCaseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "research_case_id tidak boleh kosong",
		})
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

	/* ---------- NOTIFIKASI KE DOSPEM ---------- */
	if user.Student != nil && user.Student.SupervisorID != "" {
		supervisorID := user.Student.SupervisorID

		notif := &models.Notification{
			UserID:  &supervisorID,            // target: dospem
			Type:    "weekly_report_submitted",
			Message: fmt.Sprintf(
				"%s mengirim laporan mingguan (Minggu %d)",
				user.Name, input.Week,
			),
			Metadata: util.JSONB{
				"student_id":   user.Id,
				"student_name": user.Name,
				"week":         input.Week,
				"report_status":"Menunggu Review",
			},
		}
		_ = h.notificationService.CreateNotification(notif) // abaikan error
	}
	/* ------------------------------------------ */

	return c.JSON(fiber.Map{
		"message": "Laporan berhasil dikirim",
	})
}

func (h *Handler) SubmitCompanyWeeklyReport(c *fiber.Ctx) error {
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

	input := dto.CreateCompanyWeeklyReportDTO{
		ResearchCaseID: parseString(form.Value["research_case_id"]),
		Week:           parseInt(form.Value["week"]),
		Activities:     parseString(form.Value["activities"]),
		Issues:         parseString(form.Value["issues"]),
		Hopes:          parseString(form.Value["hopes"]),
		Notes:          parseString(form.Value["notes"]),
		StartDate:      parseString(form.Value["start_date"]),
		EndDate:        parseString(form.Value["end_date"]),
	}

	// Validasi
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
	}

	// Upload file
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

	// Simpan ke service
	if err := h.weeklyReportService.SubmitCompanyReport(user.Id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan laporan perusahaan",
			"error":   err.Error(),
		})
	}

		// -- Kirim notifikasi ke perusahaan --
	researchCase, err := h.researchCaseService.GetResearchCaseByID(input.ResearchCaseID)
	if err == nil && researchCase.CompanyID != "" {
		companyID := researchCase.CompanyID

		notif := &models.Notification{
			CompanyID: &companyID,
			Type:      "student_company_report",
			Message:   fmt.Sprintf("Mahasiswa %s mengirim laporan minggu ke-%d", user.Name, input.Week),
			Metadata: util.JSONB{
				"student_id":        user.Id,
				"student_name":      user.Name,
				"research_case_id":  researchCase.ID,
				"research_case_title": researchCase.Title,
				"week":              input.Week,
				"report_type":       "company_weekly_report",
			},
		}
		_ = h.notificationService.CreateNotification(notif) // abaikan error
	}
	// -------------------------------------

	return c.JSON(fiber.Map{
		"message": "Laporan perusahaan berhasil dikirim",
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

func (h *Handler) GetWeeklyReportsForCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Pastikan user adalah perusahaan
	if user.RoleId != 2 || user.Company == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses hanya untuk user perusahaan",
		})
	}

	companyID := user.Company.CompanyID

	// Ambil laporan dari service
	reports, err := h.weeklyReportService.GetReportsByCompany(companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil laporan mingguan perusahaan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil laporan",
		"data":    reports,
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
