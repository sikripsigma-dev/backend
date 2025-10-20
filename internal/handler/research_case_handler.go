package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ResearchCaseHandler struct {
	researchCaseService *service.ResearchCaseService
}

func NewResearchCaseHandler(researchCaseService *service.ResearchCaseService) *ResearchCaseHandler {
	return &ResearchCaseHandler{researchCaseService}
}

func (h *ResearchCaseHandler) CreateResearchCase(c *fiber.Ctx) error {
	var req dto.CreateResearchCaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	researchCase, err := h.researchCaseService.CreateResearchCase(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "Research case created successfully",
		"research_case": researchCase,
	})
}

func (h *ResearchCaseHandler) GetAllResearchCases(c *fiber.Ctx) error {
	researchCases, err := h.researchCaseService.GetAllResearchCases()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"research_cases": researchCases})
}

func (h *ResearchCaseHandler) GetAllResearchCasesByStudent(c *fiber.Ctx) error {
	// studentID := c.Locals("user_id").(string) // misal dari middleware JWT
	user := c.Locals("user").(*models.User)
	studentID := user.Student.UserID

	researchCases, err := h.researchCaseService.GetApprovedByHeadStudy(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"research_cases": researchCases,
	})
}


// get all research cases for headstudy (not published yet)
func (h *ResearchCaseHandler) GetAllResearchCasesForHeadstudy(c *fiber.Ctx) error {

	// headstudyUserID := c.Locals("headstudy_user_id").(string)
	user := c.Locals("user").(*models.User)
	headstudyUserID := user.Id
	// headstudyUserID := "test-sparda-123"

	researchCases, err := h.researchCaseService.GetAllResearchCasesForHeadstudy(headstudyUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"research_cases": researchCases})
}


func (h *ResearchCaseHandler) GetResearchCaseByID(c *fiber.Ctx) error {
	id := c.Params("id")

	researchCase, err := h.researchCaseService.GetResearchCaseByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"research_case": researchCase})
}

// get research cases by company ID
func (h *ResearchCaseHandler) GetResearchCasesByCompanyID(c *fiber.Ctx) error {
	companyID := c.Params("company_id")

	researchCases, err := h.researchCaseService.GetResearchCasesByCompanyID(companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"research_cases": researchCases})
}

func (h *ResearchCaseHandler) UpdateResearchCase(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateResearchCaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	researchCase, err := h.researchCaseService.UpdateResearchCase(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":       "Research case updated successfully",
		"research_case": researchCase,
	})
}

func (h *ResearchCaseHandler) SetActiveStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	type request struct {
		IsActive bool `json:"is_active"`
	}
	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.researchCaseService.SetResearchCaseActiveStatus(id, body.IsActive); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status berhasil diperbarui",
	})
}