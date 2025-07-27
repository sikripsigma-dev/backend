package service

import (
	"Skripsigma-BE/internal/constants"
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"

	"github.com/google/uuid"
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

func (s *ChatService) CreateOrGetChatRoom(initiatorID string, targetUserID *string, targetCompanyID *string) (*models.ChatRooms, error) {
	// Cek apakah room sudah ada
	existingRoom, err := s.chatRepository.FindRoomBetween(initiatorID, targetUserID, targetCompanyID)
	if err == nil && existingRoom != nil {
		return existingRoom, nil // room sudah ada
	}

	// Room tidak ada, maka buat baru
	newRoom := &models.ChatRooms{
		ID:              uuid.NewString(),
		InitiatorID:     initiatorID,
		TargetUserID:    targetUserID,
		TargetCompanyID: targetCompanyID,
	}

	if err := s.chatRepository.CreateChatRoom(newRoom); err != nil {
		return nil, fmt.Errorf("gagal membuat chat room: %w", err)
	}

	return newRoom, nil
}


func (s *ChatService) GetChatRoomsByUserID(userID string) ([]dto.ChatRoomsWithLatestMessage, error) {
	chatRooms, err := s.chatRepository.GetRoomsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data chat rooms: %w", err)
	}

	var result []dto.ChatRoomsWithLatestMessage

	for _, room := range chatRooms {
		// ── pesan terakhir ───────────────────────────────────────────────
		lastMessage, lastMessageAt := "", ""
		if len(room.Messages) > 0 {
			lastMessage   = room.Messages[0].Message
			lastMessageAt = room.Messages[0].CreatedAt.Format("2006-01-02 15:04:05")
		}

		// ── lawan bicara  & gambar ──────────────────────────────────────
		withName, withType, withImage := "", "", ""

		/* ================  SAYA = INITIATOR  ========================== */
		if room.InitiatorID == userID {

			// ► target = perusahaan langsung
			if room.TargetCompany != nil && room.TargetCompanyID != nil {
				withName  = room.TargetCompany.Name
				withType  = "company"
				withImage = room.TargetCompany.Logo

			// ► target = user
			} else if room.TargetUser != nil && room.TargetUserID != nil {

				// user perusahaan
				if room.TargetUser.RoleId == constants.RoleCompany &&
				   room.TargetUser.Company != nil {

					withName  = room.TargetUser.Company.Company.Name
					withType  = "company"
					withImage = room.TargetUser.Company.Company.Logo

				// user biasa (student / supervisor)
				} else {
					withName  = room.TargetUser.Name
					withType  = roleLabel(room.TargetUser.RoleId)
					withImage = room.TargetUser.Image
				}
			}

		/* ================  SAYA BUKAN INITIATOR  ====================== */
		} else {

			// ► initiator = user perusahaan
			if room.Initiator.RoleId == constants.RoleCompany &&
			   room.Initiator.Company != nil {

				withName  = room.Initiator.Company.Company.Name
				withType  = "company"
				withImage = room.Initiator.Company.Company.Logo

			// ► initiator = user biasa
			} else {
				withName  = room.Initiator.Name
				withType  = roleLabel(room.Initiator.RoleId)
				withImage = room.Initiator.Image
			}
		}

		// ── append ke result ────────────────────────────────────────────
		result = append(result, dto.ChatRoomsWithLatestMessage{
			ID:              room.ID,
			InitiatorID:     room.InitiatorID,
			TargetUserID:    room.TargetUserID,
			TargetCompanyID: room.TargetCompanyID,
			LastMessage:     lastMessage,
			LastMessageAt:   lastMessageAt,
			UnreadCount:     room.UnreadCount,
			Image:           withImage,
			CreatedAt:       room.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       room.UpdatedAt.Format("2006-01-02 15:04:05"),
			WithName:        withName,
			WithType:        withType,
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

func roleLabel(roleID uint) string {
	switch roleID {
	case constants.RoleAdmin:
		return "admin"
	case constants.RoleCompany:
		return "company"
	case constants.RoleStudent:
		return "student"
	case constants.RoleSupervisor:
		return "supervisor"
	default:
		return "unknown"
	}
}


