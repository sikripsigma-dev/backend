package controllers

import (
	"Skripsigma-BE/database"
	"Skripsigma-BE/models"

	"github.com/gofiber/fiber/v2"
)

// Create Tag
func CreateTag(c *fiber.Ctx) error {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi: pastikan name tidak kosong
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tag name is required"})
	}

	// Cek apakah tag sudah ada
	var existingTag models.Tag
	if err := database.DB.Where("name = ?", input.Name).First(&existingTag).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Tag already exists"})
	}

	// Buat tag baru
	tag := models.Tag{
		Name: input.Name,
	}

	// Simpan ke database
	if err := database.DB.Create(&tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tag"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tag created successfully",
		"data":    tag,
	})
}

// Get All Tags
func GetTags(c *fiber.Ctx) error {
	var tags []models.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tags"})
	}
	return c.JSON(fiber.Map{"data": tags})
}

// Get Tag by ID
func GetTag(c *fiber.Ctx) error {
	id := c.Params("id")
	var tag models.Tag

	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tag not found"})
	}

	return c.JSON(fiber.Map{"data": tag})
}

// Update Tag
func UpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")
	var tag models.Tag

	// Cek apakah tag ada
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tag not found"})
	}

	// Parsing body request
	var input struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi: pastikan Name tidak kosong
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tag name is required"})
	}

	// Update field
	tag.Name = input.Name

	// Simpan perubahan
	if err := database.DB.Save(&tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update tag"})
	}

	return c.JSON(fiber.Map{
		"message": "Tag updated successfully",
		"data":    tag,
	})
}

// Delete Tag
func DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	var tag models.Tag

	// Cek apakah tag ada
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tag not found"})
	}

	// Hapus dari database
	if err := database.DB.Delete(&tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete tag"})
	}

	return c.JSON(fiber.Map{"message": "Tag deleted successfully"})
}
