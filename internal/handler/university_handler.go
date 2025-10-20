package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UniversityHandler struct {
	service service.UniversityService
}

func NewUniversityHandler(service service.UniversityService) *UniversityHandler {
	return &UniversityHandler{service}
}

type UniversityListResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Data    []models.University `json:"data,omitempty"`
	Error   string              `json:"error,omitempty"`
}

type UniversityResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    *models.University `json:"data,omitempty"`
	Error   string             `json:"error,omitempty"`
}

func isUUIDUniversity(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func (h *UniversityHandler) GetAll(c *fiber.Ctx) error {
	unis, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(UniversityListResponse{
			Success: false, Error: "Gagal memuat data universitas",
		})
	}
	return c.JSON(UniversityListResponse{
		Success: true, Data: unis,
	})
}

func (h *UniversityHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if !isUUIDUniversity(id) {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Parameter id bukan UUID yang valid",
		})
	}
	uni, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(UniversityResponse{
			Success: false, Error: "Universitas tidak ditemukan",
		})
	}
	return c.JSON(UniversityResponse{Success: true, Data: uni})
}

type universityRequest struct {
	Name string `json:"name"`
}

func (h *UniversityHandler) Create(c *fiber.Ctx) error {
	var req universityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Format JSON tidak valid",
		})
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Nama universitas wajib diisi",
		})
	}
	if len(req.Name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Nama universitas maksimal 255 karakter",
		})
	}

	uni, err := h.service.Create(c.Context(), req.Name)
	if err != nil {
		msg := strings.ToLower(err.Error())
		code := fiber.StatusInternalServerError
		out := err.Error()
		if strings.Contains(msg, "sudah digunakan") {
			code = fiber.StatusConflict
			out = "Universitas dengan nama tersebut sudah ada"
		}
		return c.Status(code).JSON(UniversityResponse{Success: false, Error: out})
	}

	return c.Status(fiber.StatusCreated).JSON(UniversityResponse{
		Success: true,
		Message: "Universitas berhasil dibuat",
		Data:    uni,
	})
}

func (h *UniversityHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if !isUUIDUniversity(id) {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Parameter id bukan UUID yang valid",
		})
	}

	var req universityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Format request tidak valid",
		})
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Nama universitas wajib diisi",
		})
	}
	if len(req.Name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Nama universitas maksimal 255 karakter",
		})
	}

	if err := h.service.Update(c.Context(), id, req.Name); err != nil {
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "tidak ditemukan"):
			return c.Status(fiber.StatusNotFound).JSON(UniversityResponse{
				Success: false, Error: err.Error(),
			})
		case strings.Contains(msg, "sudah digunakan"):
			return c.Status(fiber.StatusConflict).JSON(UniversityResponse{
				Success: false, Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(UniversityResponse{
				Success: false, Error: "Terjadi kesalahan internal",
			})
		}
	}

	return c.JSON(UniversityResponse{
		Success: true, Message: "Universitas berhasil diupdate",
	})
}

func (h *UniversityHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if !isUUIDUniversity(id) {
		return c.Status(fiber.StatusBadRequest).JSON(UniversityResponse{
			Success: false, Error: "Parameter id bukan UUID yang valid",
		})
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		// Service sudah mengubah not found menjadi pesan yang jelas
		// Di sini kita anggap itu 404; selain itu 500.
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(UniversityResponse{
				Success: false, Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(UniversityResponse{
			Success: false, Error: "Terjadi kesalahan internal",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
