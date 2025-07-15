package handler

import (
	"Skripsigma-BE/internal/dto"
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