package controllers

import (
	"Skripsigma-BE/database"
	"Skripsigma-BE/models"

	"github.com/gofiber/fiber/v2"
)

// CreateCompany - Menambahkan perusahaan baru
func CreateCompany(c *fiber.Ctx) error {
	var company models.Company

	// Parse request body ke struct company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Simpan ke database
	if err := database.DB.Create(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create company"})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Company created successfully",
		"company": company,
	})
}

// GetAllCompanies - Mendapatkan semua perusahaan
func GetAllCompanies(c *fiber.Ctx) error {
	var companies []models.Company

	// Ambil semua data perusahaan
	if err := database.DB.Find(&companies).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve companies"})
	}

	// Return response
	return c.JSON(companies)
}

// GetCompanyByID - Mendapatkan perusahaan berdasarkan ID
func GetCompanyByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company

	// Cari perusahaan berdasarkan ID
	if err := database.DB.First(&company, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}

	// Return response
	return c.JSON(company)
}

// UpdateCompany - Mengupdate data perusahaan berdasarkan ID
func UpdateCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company

	// Cek apakah perusahaan ada
	if err := database.DB.First(&company, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}

	// Parse data baru
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update data perusahaan
	if err := database.DB.Save(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update company"})
	}

	// Return response
	return c.JSON(fiber.Map{"message": "Company updated successfully", "company": company})
}

// DeleteCompany - Menghapus perusahaan berdasarkan ID
func DeleteCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company

	// Cek apakah perusahaan ada
	if err := database.DB.First(&company, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}

	// Hapus perusahaan dari database
	if err := database.DB.Delete(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete company"})
	}

	// Return response
	return c.JSON(fiber.Map{"message": "Company deleted successfully"})
}