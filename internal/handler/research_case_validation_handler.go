package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"strings"

	"github.com/gofiber/fiber/v2"
)

type ResearchCaseValidationHandler struct {
	researchCaseValidationService *service.ResearchCaseValidationService
}

func NewResearchCaseValidationHandler(researchCaseValidationService *service.ResearchCaseValidationService) *ResearchCaseValidationHandler {
	return &ResearchCaseValidationHandler{researchCaseValidationService}
}

func (h *ResearchCaseValidationHandler) CreateValidation(c *fiber.Ctx) error {
	var req dto.CreateResearchCaseValidationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validation, err := h.researchCaseValidationService.CreateValidation(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Validation created successfully",
		"validation": validation,
	})
}

func (h *ResearchCaseValidationHandler) GetValidationsByResearchCaseID(c *fiber.Ctx) error {
	researchCaseID := c.Params("research_case_id")

	validations, err := h.researchCaseValidationService.GetValidationsByResearchCaseID(researchCaseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"validations": validations})
}

// func (h *ResearchCaseValidationHandler) CreateComment(c *fiber.Ctx) error {
// 	var req dto.CreateResearchCaseValidationCommentRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	comment, err := h.researchCaseValidationService.CreateComment(req)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Comment created successfully",
// 		"comment": comment,
// 	})
// }

// handler/research_case_validation_handler.go
func (h *ResearchCaseValidationHandler) CreateComment(c *fiber.Ctx) error {
    // ambil id study case dari path
    rcID := c.Params("research_case_id")

    var req dto.CreateResearchCaseValidationCommentRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // pastikan research_case_id terisi (prioritas dari path param)
    if rcID == "" && strings.TrimSpace(req.ResearchCaseID) == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "research_case_id required"})
    }
    if rcID != "" {
        req.ResearchCaseID = rcID
    }

    // ambil user login dari context lalu pakai sebagai program_head_id
    u, ok := c.Locals("user").(*models.User)
    if !ok || u == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }
    req.ProgramHeadID = u.Id // ← penting, supaya FK valid

    if strings.TrimSpace(req.Comment) == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "comment required"})
    }

    comment, err := h.researchCaseValidationService.CreateComment(req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"comment": comment})
}


func (h *ResearchCaseValidationHandler) GetCommentsByResearchCaseID(c *fiber.Ctx) error {
	researchCaseID := c.Params("research_case_id")

	comments, err := h.researchCaseValidationService.GetCommentsByResearchCaseID(researchCaseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"comments": comments})
}

func (h *ResearchCaseValidationHandler) UpsertValidation(c *fiber.Ctx) error {
    var req dto.UpdateResearchCaseValidationRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    rcID := c.Query("research_case_id")
    if rcID == "" {
        rcID = strings.TrimSpace(req.ResearchCaseID) // fallback dari body
    }
    if rcID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "research_case_id required"})
    }

    user := c.Locals("user").(*models.User) // kaprodi login
    if err := h.researchCaseValidationService.UpsertValidationStatus(user.Id, rcID, req.Status); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Status updated"})
}