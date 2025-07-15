package handler

import (
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


func (h *ChatHandler) GetChatRoomsByStudentID(c *fiber.Ctx) error {
	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "student ID is required",
		})
	}

	chatRooms, err := h.chatService.GetChatRoomsByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": chatRooms,
	})
}

func (h *ChatHandler) GetChatRoomByCompanyID(c *fiber.Ctx) error{
	companyID := c.Params("id")
	if companyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "student ID is required",
		})
	}

	chatRooms, err := h.chatService.GetChatRoomsByCompanyID(companyID)
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
