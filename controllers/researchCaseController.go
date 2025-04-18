package controllers

import (
	"Skripsigma-BE/database"
	"Skripsigma-BE/models"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Create Research Case
func CreateResearchCase(c *fiber.Ctx) error {
	// Parsing body request
	var input struct {
		CompanyID            string   `json:"company_id"`
		Title                string   `json:"title"`
		Field                string   `json:"field"`
		EducationRequirement string   `json:"education_requirement"`
		Duration             string   `json:"duration"`
		Description          string   `json:"description"`
		TagIDs               []string `json:"tag_ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi: Pastikan CompanyID tidak kosong
	if input.CompanyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Company ID is required"})
	}

	// Cek apakah perusahaan dengan CompanyID ini ada di database
	var company models.Company
	if err := database.DB.Where("id = ?", input.CompanyID).First(&company).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}

	// Buat research case baru
	researchCase := models.ResearchCase{
		CompanyID:            input.CompanyID,
		Title:                input.Title,
		Field:                input.Field,
		EducationRequirement: input.EducationRequirement,
		Duration:             input.Duration,
		Description:          input.Description,
	}

	// Simpan ke database
	if err := database.DB.Create(&researchCase).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create research case"})
	}

	// Cek apakah ada tag yang dikirim
	if len(input.TagIDs) > 0 {
		for _, tagID := range input.TagIDs {
			// Debugging: pastikan ID valid
			fmt.Println("Inserting tag:", tagID)

			// Masukkan ke table ss_t_research_case_tags
			researchCaseTag := models.ResearchCaseTag{
				ResearchCaseID: researchCase.ID,
				TagID:          tagID,
			}
			if err := database.DB.Create(&researchCaseTag).Error; err != nil {
				fmt.Println("Error inserting tag:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to associate tags with research case",
				})
			}
		}
	}

	// Ambil research case yang baru dibuat, beserta relasi company dan tags
	var createdResearchCase models.ResearchCase
	if err := database.DB.Preload("Company").Preload("Tags").First(&createdResearchCase, "id = ?", researchCase.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve created research case"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Research case created successfully",
		"data":    researchCase,
	})
}

// Get All Research Cases
func GetResearchCases(c *fiber.Ctx) error {
	var researchCases []models.ResearchCase
	if err := database.DB.Preload("Company").Find(&researchCases).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch research cases"})
	}
	return c.JSON(fiber.Map{"data": researchCases})
}

// Get Research Case by ID
func GetResearchCase(c *fiber.Ctx) error {
	id := c.Params("id") 
	var researchCase models.ResearchCase
	log.Println("Fetching research case with ID:", id)
	// Preload Company dan Tags agar data terkait ikut terambil
	if err := database.DB.Preload("Company").Preload("Tags").First(&researchCase, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Research case not found"})
	}

	return c.JSON(fiber.Map{"data": researchCase})
}

// Update Research Case
func UpdateResearchCase(c *fiber.Ctx) error {
	id := c.Params("id")
	var researchCase models.ResearchCase

	// Cek apakah Research Case ada
	if err := database.DB.First(&researchCase, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Research case not found"})
	}

	// Parsing body request
	var input struct {
		CompanyID            string `json:"company_id"`
		Title                string `json:"title"`
		Field                string `json:"field"`
		EducationRequirement string `json:"education_requirement"`
		Duration             string `json:"duration"`
		Description          string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi: Pastikan CompanyID ada di database jika diubah
	if input.CompanyID != "" {
		var company models.Company
		if err := database.DB.First(&company, "id = ?", input.CompanyID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
		}
	}

	// Update field jika tidak kosong
	if input.CompanyID != "" {
		researchCase.CompanyID = input.CompanyID
	}
	if input.Title != "" {
		researchCase.Title = input.Title
	}
	if input.Field != "" {
		researchCase.Field = input.Field
	}
	if input.EducationRequirement != "" {
		researchCase.EducationRequirement = input.EducationRequirement
	}
	if input.Duration != "" {
		researchCase.Duration = input.Duration
	}
	if input.Description != "" {
		researchCase.Description = input.Description
	}

	// Simpan perubahan
	if err := database.DB.Save(&researchCase).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update research case"})
	}

	return c.JSON(fiber.Map{
		"message": "Research case updated successfully",
		"data":    researchCase,
	})
}

// Delete Research Case
func DeleteResearchCase(c *fiber.Ctx) error {
	id := c.Params("id")
	var researchCase models.ResearchCase

	// Cek apakah Research Case ada
	if err := database.DB.First(&researchCase, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Research case not found"})
	}

	// Hapus dari database
	if err := database.DB.Delete(&researchCase).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete research case"})
	}

	return c.JSON(fiber.Map{"message": "Research case deleted successfully"})
}