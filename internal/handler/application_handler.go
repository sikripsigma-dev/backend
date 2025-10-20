package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"Skripsigma-BE/internal/util"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)


type ApplicationHandler struct {
	applicationService  *service.ApplicationService
	notificationService service.NotificationService
	researchCaseService *service.ResearchCaseService
}

func NewApplicationHandler(
	applicationService *service.ApplicationService,
	notificationService service.NotificationService,
	researchCaseService *service.ResearchCaseService,
) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService:  applicationService,
		notificationService: notificationService,
		researchCaseService: researchCaseService,
	}
}


func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	// 1. Parse JSON body ke struct
	var req dto.CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("[ERROR] BodyParser:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
	}

	// 2. Validasi manual jika perlu
	if req.ResearchCaseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ResearchCaseID wajib diisi",
		})
	}

	// 3. Ambil user dari context
	user := c.Locals("user").(*models.User)

	// 4. Simpan aplikasi melalui service
	application, err := h.applicationService.CreateApplication(req, user.Id)
	if err != nil {
		log.Println("[ERROR] CreateApplication:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan aplikasi",
			"error":   err.Error(),
		})
	}

	// 5. (Opsional) Kirim notifikasi ke pemilik studi kasus
	researchCase, err := h.researchCaseService.GetResearchCaseByID(req.ResearchCaseID)
	if err == nil && researchCase.CompanyID != "" {
		companyID := researchCase.CompanyID
		caseTitle := researchCase.Title

		notif := &models.Notification{
			CompanyID: &companyID,
			Type:      "apply",
			Message:   fmt.Sprintf("%s mengajukan studi kasus '%s'", user.Name, caseTitle),
			Metadata: util.JSONB{
				"case_id":    req.ResearchCaseID,
				"case_title": caseTitle,
				"applicant":  user.Name,
			},
		}

		_ = h.notificationService.CreateNotification(notif) // abaikan error
	}

	// 6. Return hasil sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil mengajukan studi kasus",
		"application": dto.ApplyResponse{
			ID:             application.ID,
			ResearchCaseID: application.ResearchCaseID,
			UserID:         application.UserID,
			Status:         application.Status,
			AppliedAt:      application.AppliedAt,
			ProcessedAt:    application.ProcessedAt,
		},
	})
}


// cek application exist
func (h *ApplicationHandler) CheckApplication(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	researchCaseID := c.Query("research_case_id")

	assigned, err := h.applicationService.CheckStudentAlreadyAssigned(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengecek status penugasan",
		})
	}

	exists, err := h.applicationService.CheckApplicationExists(user.Id, researchCaseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengecek aplikasi",
		})
	}

	return c.JSON(fiber.Map{
		"applied":  exists,
		"assigned": assigned,
		"eligible": !exists && !assigned,
	})
}

func (h *ApplicationHandler) ProcessApplication(c *fiber.Ctx) error {
	// 1. Parse body request
	var req dto.ProcessApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	// 2. Ambil ID aplikasi dan user login
	applicationID := c.Params("id")
	user := c.Locals("user").(*models.User)

	// 3. Proses aplikasi
	application, err := h.applicationService.ProcessApplication(applicationID, req.Status, user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Ambil data studi kasus untuk notifikasi (judul)
	researchCase, err := h.researchCaseService.GetResearchCaseByID(application.ResearchCaseID)
	caseTitle := "(judul tidak ditemukan)"
	if err == nil && researchCase.Title != "" {
		caseTitle = researchCase.Title
	}

	// 5. Siapkan notifikasi ke mahasiswa
	statusLabel := map[string]string{
		"accepted":      "diterima",
		"rejected":      "ditolak",
		"need_revision": "memerlukan revisi",
	}[strings.ToLower(application.Status)]
	if statusLabel == "" {
		statusLabel = application.Status // fallback
	}

	notif := &models.Notification{
		UserID:  &application.UserID,
		Type:    "application_review",
		Message: fmt.Sprintf("Aplikasi Anda untuk studi kasus '%s' telah %s", caseTitle, statusLabel),
		Metadata: util.JSONB{
			"application_id": application.ID,
			"case_id":        application.ResearchCaseID,
			"case_title":     caseTitle,
			"status":         application.Status,
			"reviewer_id":    user.Id,
		},
	}

	// 6. Simpan notifikasi (error tidak memblokir)
	_ = h.notificationService.CreateNotification(notif)

	// 7. Kirim response sukses
	return c.JSON(fiber.Map{
		"message":     "Application processed successfully",
		"application": application, // atau mapping ke DTO jika ingin lebih aman
	})
}

func (h *ApplicationHandler) RespondToApplication(c *fiber.Ctx) error {
	// 1. Ambil dan parse body
	var req dto.ProcessApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 2. Param ID dan user login (mahasiswa)
	applicationID := c.Params("id")
	user := c.Locals("user").(*models.User)

	// 3. Proses via service (ubah status jadi "confirmed")
	application, err := h.applicationService.RespondToApplication(applicationID, req.Status, user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Ambil info studi kasus untuk ambil perusahaan + judul
	researchCase, err := h.researchCaseService.GetResearchCaseByID(application.ResearchCaseID)
	if err == nil && researchCase.CompanyID != "" {
		companyID := researchCase.CompanyID
		caseTitle := researchCase.Title

		// 5. Kirim notifikasi ke perusahaan
		notif := &models.Notification{
			CompanyID: &companyID,
			Type:      "application_confirmed", // sesuaikan dengan enum/type kamu
			Message:   fmt.Sprintf("%s telah mengonfirmasi penugasan untuk studi kasus '%s'", user.Name, caseTitle),
			Metadata: util.JSONB{
				"application_id": application.ID,
				"case_id":        researchCase.ID,
				"case_title":     caseTitle,
				"student_name":   user.Name,
			},
		}
		_ = h.notificationService.CreateNotification(notif) // abaikan error
	}

	// 6. Response ke frontend
	return c.JSON(fiber.Map{
		"message":     "Application responded successfully",
		"application": application,
	})
}

// get by research case ID
func (h *ApplicationHandler) GetApplicationsByResearchCaseID(c *fiber.Ctx) error {
    researchCaseID := c.Params("id")

    applications, err := h.applicationService.GetApplicationsByResearchCaseID(researchCaseID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": fmt.Sprintf("Gagal ambil data pelamar: %v", err),
        })
    }

    var responses []dto.ApplicationResponse
    for _, app := range applications {
        // Filter: hanya masukkan jika HeadStudyDecision bernilai "accepted"
        if app.HeadStudyDecision != "accepted" {
            continue
        }

        resp := dto.ApplicationResponse{
            ID:             app.ID,
            ResearchCaseID: app.ResearchCaseID,
            Status:         app.Status,
            AppliedAt:      app.AppliedAt,
            ProcessedAt:    app.ProcessedAt,
            ProcessedBy:    app.ProcessedBy,
            User: dto.ApplicationUserResponse{
                Id:    app.User.Id,
                Nim:   app.User.Nim,
                Name:  app.User.Name,
                Phone: app.User.Phone,
                Email: app.User.Email,
            },
        }
        responses = append(responses, resp)
    }

    return c.JSON(fiber.Map{
        "applications": responses,
    })
}

// get students applications for Headstudy
func (h *ApplicationHandler) GetPendingApplicationsForHeadStudy(c *fiber.Ctx) error {

	headStudyUser := c.Locals("user").(*models.User)

	if headStudyUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	applications, err := h.applicationService.GetPendingHeadStudyConfirmations(
		headStudyUser.Headstudy.UniversityID,
		headStudyUser.Headstudy.StudyProgramID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Gagal mengambil aplikasi: %v", err),
		})
	}

	var responses []dto.ApplicationResponse
	for _, app := range applications {
		

		resp := dto.ApplicationResponse{
			ID:             app.ID,
			ResearchCaseID: app.ResearchCaseID,
			Status:         app.Status,
			AppliedAt:      app.AppliedAt,
			ProcessedAt:    app.ProcessedAt,
			ProcessedBy:    app.ProcessedBy,
			User: dto.ApplicationUserResponse{
				Id:    app.User.Id,
				Nim:   app.User.Nim,
				Name:  app.User.Name,
				Phone: app.User.Phone,
				Email: app.User.Email,
			},
			ResearchCase: dto.ResearchCaseResponse{
				ID:          app.ResearchCase.ID,
				Title:       app.ResearchCase.Title,
				Description: app.ResearchCase.Description,
				Company: dto.CompanyResponse{
					ID:    app.ResearchCase.Company.Id,
					Name:  app.ResearchCase.Company.Name,
					Email: app.ResearchCase.Company.Email,
				},
			},
		}


		responses = append(responses, resp)
	}

	return c.JSON(fiber.Map{
		"applications": responses,
	})
}


type ConfirmRequest struct {
	Decision string `json:"decision"` // expected: accepted / rejected
}

func (h *ApplicationHandler) ConfirmByHeadStudy(c *fiber.Ctx) error {
	idParam := c.Params("id")
	applicationID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid application id",
		})
	}

	var req ConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	err = h.applicationService.ConfirmByHeadStudy(uint(applicationID), req.Decision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "decision updated successfully",
	})
}


func (h *ApplicationHandler) GetApplicationsByStudentID(c *fiber.Ctx) error {
	studentID := c.Params("id")

	applications, err := h.applicationService.GetApplicationsByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var responses []dto.GetApplicationByStudentResponse
	for _, app := range applications {
		// Map tags
		var tags []dto.TagResponse
		for _, tag := range app.ResearchCase.Tags {
			tags = append(tags, dto.TagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}

		// Map response
		res := dto.GetApplicationByStudentResponse{
			ID:          app.ID,
			Status:      app.Status,
			AppliedAt:   app.AppliedAt,
			ProcessedAt: &app.ProcessedAt,
			ProcessedBy: app.ProcessedBy,
			ResearchCase: dto.ResearchCaseResponse{
				ID:                   app.ResearchCase.ID,
				Title:                app.ResearchCase.Title,
				Field:                app.ResearchCase.Field,
				Location:             app.ResearchCase.Location,
				EducationRequirement: app.ResearchCase.EducationRequirement,
				Duration:             app.ResearchCase.Duration,
				Description:          app.ResearchCase.Description,
				CreatedAt:            app.ResearchCase.CreatedAt,
				Tags:                 tags,
				Company: dto.CompanyResponse{
					ID:          app.ResearchCase.Company.Id,
					Name:        app.ResearchCase.Company.Name,
					Email:       app.ResearchCase.Company.Email,
					Phone:       app.ResearchCase.Company.Phone,
					Address:     app.ResearchCase.Company.Address,
					Description: app.ResearchCase.Company.Description,
				},
			},
		}

		responses = append(responses, res)
	}

	return c.JSON(fiber.Map{
		"message":      "Applications fetched successfully",
		"applications": responses,
	})
}

