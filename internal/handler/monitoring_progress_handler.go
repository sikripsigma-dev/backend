package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MonitoringProgressHandler struct {
	service service.MonitoringProgressService
}

func NewMonitoringProgressHandler(s service.MonitoringProgressService) *MonitoringProgressHandler {
	return &MonitoringProgressHandler{
		service: s,
	}
}

func (h *MonitoringProgressHandler) GiveFeedback(c *fiber.Ctx) error {
	var payload dto.CreateSupervisorFeedbackDTO

	// Binding dari body request (JSON/form)
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
	}

	// Validasi manual kalau perlu (opsional, jika tag binding tidak cukup)
	if payload.WeeklyReportID == 0 || payload.Feedback == "" || payload.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Semua field wajib diisi",
		})
	}

	if err := h.service.GiveFeedback(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan feedback",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Feedback berhasil disimpan",
	})
}

func (h *MonitoringProgressHandler) GiveFeedbackCompany(c *fiber.Ctx) error {
	var payload dto.CreateCompanyFeedbackDTO

	// Binding dari body request (JSON/form)
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
	}

	// Validasi manual kalau perlu (opsional, jika tag binding tidak cukup)
	if payload.WeeklyReportID == 0 || payload.Feedback == ""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Semua field wajib diisi",
		})
	}

	if err := h.service.GiveFeedbackCompany(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan feedback",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Feedback berhasil disimpan",
	})
}