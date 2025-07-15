package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)


type CompanyHandler struct {
	companyService *service.CompanyService
}

func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService}
}

func (h *CompanyHandler) CreateCompany(c *fiber.Ctx) error {
	var req dto.CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	company, err := h.companyService.CreateCompany(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Company created successfully",
		"company": company,
	})
}

func (h *CompanyHandler) GetCompanyByID(c *fiber.Ctx) error {
	id := c.Params("id")
	company, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}

	return c.JSON(fiber.Map{"company": company})
}

func (h *CompanyHandler) GetAllCompanies(c *fiber.Ctx) error {
	companies, err := h.companyService.GetAllCompanies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"companies": companies})
}

func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	company, err := h.companyService.UpdateCompany(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":  "Company updated successfully",
		"company": company,
	})
}