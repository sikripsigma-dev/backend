package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"Skripsigma-BE/internal/util"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService}
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	companyID := ""
	if user.Company != nil {
		companyID = user.Company.CompanyID
	}

	notifs, err := h.notificationService.GetForReceiver(user.Id, companyID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil notifikasi"})
	}

	return c.JSON(fiber.Map{"notifications": notifs})
}

// PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*models.User)

	companyID := ""
	if user.Company != nil {
		companyID = user.Company.CompanyID
	}

	notif, err := h.notificationService.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notifikasi tidak ditemukan"})
	}
	
	if !isNotificationOwner(user.Id, companyID, notif) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Tidak diizinkan mengakses notifikasi ini",
		})
	}
	

	if err := h.notificationService.MarkAsRead(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menandai notifikasi sebagai dibaca"})
	}

	return c.JSON(fiber.Map{"message": "Notifikasi dibaca"})
}


// PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	companyID := ""
	if user.Company != nil {
		companyID = user.Company.CompanyID
	}

	if err := h.notificationService.MarkAllAsRead(user.Id, companyID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menandai semua notifikasi sebagai dibaca"})
	}

	return c.JSON(fiber.Map{"message": "Semua notifikasi dibaca"})
}

func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	// user := c.Locals("user").(*models.User)

	var payload struct {
		UserID    *string        `json:"user_id"`
		CompanyID *string        `json:"company_id"`
		Type      string         `json:"type"`
		Message   string         `json:"message"`
		Metadata  util.JSONB     `json:"metadata"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	if payload.Type == "" || payload.Message == "" || (payload.UserID == nil && payload.CompanyID == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Field wajib tidak boleh kosong"})
	}

	notif := &models.Notification{
		UserID:    payload.UserID,
		CompanyID: payload.CompanyID,
		Type:      payload.Type,
		Message:   payload.Message,
		Metadata:  payload.Metadata,
	}

	if err := h.notificationService.CreateNotification(notif); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat notifikasi"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Notifikasi berhasil dibuat"})
}


func isNotificationOwner(currentUserID, currentCompanyID string, notif *models.Notification) bool {
	if notif.UserID != nil && *notif.UserID == currentUserID {
		return true
	}
	if notif.CompanyID != nil && *notif.CompanyID == currentCompanyID {
		return true
	}
	return false
}