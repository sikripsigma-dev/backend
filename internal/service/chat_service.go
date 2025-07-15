package service

import (
	"Skripsigma-BE/internal/dto"
	// "Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
	// "gorm.io/gorm"
)

type ChatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) *ChatService {
	return &ChatService{
		chatRepository: chatRepo,
	}
}

func (s *ChatService) GetChatRoomsByStudentID(studentID string) ([]dto.ChatRoomsWithLatestMessage, error) {
	chatRooms, err := s.chatRepository.GetChatRoomsByStudentID(studentID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data chat rooms: %w", err)
	}

	var result []dto.ChatRoomsWithLatestMessage
	for _, room := range chatRooms {
		var lastMessage string
		var lastMessageAt string
		if len(room.Messages) > 0 {
			lastMessage = room.Messages[0].Message
			lastMessageAt = room.Messages[0].CreatedAt.Format("2006-01-02 15:04:05")
		}

		result = append(result, dto.ChatRoomsWithLatestMessage{
			ID:            room.ID,
			StudentID:     room.StudentID,
			CompanyID:     room.CompanyID,
			CompanyName:   room.Company.Name,
			LastMessage:   lastMessage,
			LastMessageAt: lastMessageAt,
			UnreadCount:   room.UnreadCount,
			CreatedAt:     room.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     room.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

func (s *ChatService) GetChatRoomsByCompanyID(companyID string) ([]dto.ChatRoomsWithLatestMessage, error) {
	chatRooms, err := s.chatRepository.GetChatRoomsByCompanyID(companyID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data chat rooms: %w", err)
	}

	var result []dto.ChatRoomsWithLatestMessage
	for _, room := range chatRooms {
		var lastMessage string
		var lastMessageAt string
		if len(room.Messages) > 0 {
			lastMessage = room.Messages[0].Message
			lastMessageAt = room.Messages[0].CreatedAt.Format("2006-01-02 15:04:05")
		}

		result = append(result, dto.ChatRoomsWithLatestMessage{
			ID:            room.ID,
			StudentID:     room.StudentID,
			CompanyID:     room.CompanyID,
			CompanyName:   room.Company.Name,
			LastMessage:   lastMessage,
			LastMessageAt: lastMessageAt,
			UnreadCount:   room.UnreadCount,
			CreatedAt:     room.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     room.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

func (s *ChatService) GetMessagesByRoomID(roomID string) ([]dto.ChatMessageResponse, error) {
	messages, err := s.chatRepository.GetMessagesByRoomID(roomID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil chat messages: %w", err)
	}

	var result []dto.ChatMessageResponse
	for _, msg := range messages {
		result = append(result, dto.ChatMessageResponse{
			ID:         msg.ID,
			RoomID:     msg.RoomID,
			SenderID:   msg.SenderID,
			SenderType: msg.SenderType,
			Message:    msg.Message,
			IsRead:     msg.IsRead,
			CreatedAt:  msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}


