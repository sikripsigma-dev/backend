package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StudentDocumentHandler struct {
	studentDocumentService service.StudentDocumentService
}

func NewStudentDocumentHandler(s service.StudentDocumentService) *StudentDocumentHandler {
	return &StudentDocumentHandler{studentDocumentService: s}
}

// POST /api/student-documents/upload
func (h *StudentDocumentHandler) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	file, err := c.FormFile("document")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File dokumen tidak ditemukan"})
	}

	docType := c.FormValue("type")
	if docType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tipe dokumen wajib diisi"})
	}

	document, err := h.studentDocumentService.Create(user.Id, docType, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(document)
}

func (h *StudentDocumentHandler) GetMyDocuments(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userID := user.Id

	docs, err := h.studentDocumentService.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": docs,
	})
}

func (h *StudentDocumentHandler) GetStudentDocumentByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	docs, err := h.studentDocumentService.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": docs,
	})
}