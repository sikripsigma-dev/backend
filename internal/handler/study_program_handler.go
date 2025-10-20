// internal/handler/study_program_handler.go
package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudyProgramHandler struct {
	service service.StudyProgramService
}

func NewStudyProgramHandler(service service.StudyProgramService) *StudyProgramHandler {
	return &StudyProgramHandler{service}
}

type StudyProgramListResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Data    []models.StudyProgram `json:"data,omitempty"`
	Error   string                `json:"error,omitempty"`
}

type StudyProgramResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message,omitempty"`
	Data    *models.StudyProgram `json:"data,omitempty"`
	Error   string               `json:"error,omitempty"`
}

func isUUIDStudyProgram(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

type studyProgramRequest struct {
	Name string `json:"name"`
}

func (h *StudyProgramHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(StudyProgramListResponse{
			Success: false, Error: "Gagal memuat data program studi",
		})
	}
	return c.JSON(StudyProgramListResponse{Success: true, Data: items})
}

func (h *StudyProgramHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if !isUUIDStudyProgram(id) {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Parameter id bukan UUID yang valid",
		})
	}
	item, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(StudyProgramResponse{
			Success: false, Error: "Program studi tidak ditemukan",
		})
	}
	return c.JSON(StudyProgramResponse{Success: true, Data: item})
}

func (h *StudyProgramHandler) Create(c *fiber.Ctx) error {
	var req studyProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Format JSON tidak valid",
		})
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Nama program studi wajib diisi",
		})
	}
	if len(req.Name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Nama program studi maksimal 255 karakter",
		})
	}

	item, err := h.service.Create(c.Context(), req.Name)
	if err != nil {
		msg := strings.ToLower(err.Error())
		code := fiber.StatusInternalServerError
		out := err.Error()
		if strings.Contains(msg, "sudah digunakan") {
			code = fiber.StatusConflict
			out = "Program studi dengan nama tersebut sudah ada"
		}
		return c.Status(code).JSON(StudyProgramResponse{Success: false, Error: out})
	}

	return c.Status(fiber.StatusCreated).JSON(StudyProgramResponse{
		Success: true,
		Message: "Program studi berhasil dibuat",
		Data:    item,
	})
}

func (h *StudyProgramHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if !isUUIDStudyProgram(id) {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Parameter id bukan UUID yang valid",
		})
	}
	var req studyProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Format request tidak valid",
		})
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Nama program studi wajib diisi",
		})
	}
	if len(req.Name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Nama program studi maksimal 255 karakter",
		})
	}

	if err := h.service.Update(c.Context(), id, req.Name); err != nil {
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "tidak ditemukan"):
			return c.Status(fiber.StatusNotFound).JSON(StudyProgramResponse{
				Success: false, Error: err.Error(),
			})
		case strings.Contains(msg, "sudah digunakan"):
			return c.Status(fiber.StatusConflict).JSON(StudyProgramResponse{
				Success: false, Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(StudyProgramResponse{
				Success: false, Error: "Terjadi kesalahan internal",
			})
		}
	}

	return c.JSON(StudyProgramResponse{
		Success: true, Message: "Program studi berhasil diupdate",
	})
}

func (h *StudyProgramHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if !isUUIDStudyProgram(id) {
		return c.Status(fiber.StatusBadRequest).JSON(StudyProgramResponse{
			Success: false, Error: "Parameter id bukan UUID yang valid",
		})
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(StudyProgramResponse{
				Success: false, Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(StudyProgramResponse{
			Success: false, Error: "Terjadi kesalahan internal",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
