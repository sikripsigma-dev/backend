package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"Skripsigma-BE/internal/util"
	"fmt"
	"log"

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

// func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
// 	var req dto.CreateApplicationRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	user := c.Locals("user").(*models.User)
// 	application, err := h.applicationService.CreateApplication(req, user.Id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Applied successfully",
// 		"application": dto.ApplyResponse{
// 			ID:             application.ID,
// 			ResearchCaseID: application.ResearchCaseID,
// 			UserID:         application.UserID,
// 			Status:         application.Status,
// 			AppliedAt:      application.AppliedAt,
// 			ProcessedAt:    application.ProcessedAt,
// 		},
// 	})	
// }

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
func (h *ApplicationHandler) CheckApplicationExists(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	researchCaseID := c.Query("research_case_id")

	exists, err := h.applicationService.CheckApplicationExists(user.Id, researchCaseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengecek aplikasi",
		})
	}

	return c.JSON(fiber.Map{
		"applied": exists,
	})
}

func (h *ApplicationHandler) ProcessApplication(c *fiber.Ctx) error {
	var req dto.ProcessApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	applicationID := c.Params("id")
	user := c.Locals("user").(*models.User)

	err := h.applicationService.ProcessApplication(applicationID, req.Status, user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Application processed successfully",
	})
}

func (h *ApplicationHandler) RespondToApplication(c *fiber.Ctx) error {
	var req dto.ProcessApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	applicationID := c.Params("id")
	user := c.Locals("user").(*models.User)

	err := h.applicationService.RespondToApplication(applicationID, req.Status, user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Application responded successfully",
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

