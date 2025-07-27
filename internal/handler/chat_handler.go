package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler (chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}


func (h *ChatHandler) CreateChatRoom(c *fiber.Ctx) error {
	type CreateRoomRequest struct {
		TargetUserID    *string `json:"target_user_id"`
		TargetCompanyID *string `json:"target_company_id"`
	}

	var body CreateRoomRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user := c.Locals("user").(*models.User)
	initiatorID := user.Id

	if body.TargetUserID == nil && body.TargetCompanyID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "target_user_id atau target_company_id harus diisi",
		})
	}

	room, err := h.chatService.CreateOrGetChatRoom(initiatorID, body.TargetUserID, body.TargetCompanyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": room,
	})
}

func (h *ChatHandler) GetChatRoomsByUserID(c *fiber.Ctx) error {
	// userID := c.Params("id")
	user := c.Locals("user").(*models.User)
	userID := user.Id

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	chatRooms, err := h.chatService.GetChatRoomsByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": chatRooms,
	})
}

func (h *ChatHandler) GetMessagesByRoomID(c *fiber.Ctx) error {
	roomID := c.Params("room_id")
	if roomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "room_id tidak boleh kosong",
		})
	}

	messages, err := h.chatService.GetMessagesByRoomID(roomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": messages,
	})
}
