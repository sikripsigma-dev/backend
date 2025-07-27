package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateChatRoom(room *models.ChatRooms) error
	GetRoomsByUserID(userID string) ([]models.ChatRooms, error)
	GetMessagesByRoomID(roomID string) ([]models.ChatMessage, error)
	FindRoomBetween(userID1 string, userID2 *string, companyID *string) (*models.ChatRooms, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db}
}

func (r *chatRepository) CreateChatRoom(room *models.ChatRooms) error {
	return r.db.Create(room).Error
}

func (r *chatRepository) GetRoomsByUserID(userID string) ([]models.ChatRooms, error) {
	var rooms []models.ChatRooms
	r.db.
		// Preload("Initiator").
		// Preload("TargetUser").
		Preload("Initiator.Company.Company").
  		Preload("TargetUser.Company.Company").
		Preload("TargetCompany").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(1)
		}).
		Where("initiator_id = ? OR target_user_id = ?", userID, userID).
		Or("target_company_id IN (SELECT company_id FROM ss_company_user WHERE user_id = ?)", userID).
		Find(&rooms)

	for i := range rooms {
		var count int64
		r.db.Model(&models.ChatMessage{}).
			Where("room_id = ? AND is_read = ? AND sender_id != ?", rooms[i].ID, false, userID).
			Count(&count)
		rooms[i].UnreadCount = int(count)
	}
	return rooms, nil
}

func (r *chatRepository) GetMessagesByRoomID(roomID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.
		Where("room_id = ?", roomID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) FindRoomBetween(initiatorID string, targetUserID *string, targetCompanyID *string) (*models.ChatRooms, error) {
	var room models.ChatRooms
	query := r.db.Model(&models.ChatRooms{}).
		Where("initiator_id = ?", initiatorID)

	if targetUserID != nil {
		query = query.Where("target_user_id = ? AND target_company_id IS NULL", *targetUserID)
	} else if targetCompanyID != nil {
		query = query.Where("target_company_id = ? AND target_user_id IS NULL", *targetCompanyID)
	} else {
		return nil, nil
	}

	err := query.First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func SaveChatMessage(db *gorm.DB, msg *models.ChatMessage) error {
	return db.Create(msg).Error
}

func MarkMessagesAsRead(db *gorm.DB, roomID string, readerID string) error {
	return db.Model(&models.ChatMessage{}).
		Where("room_id = ? AND is_read = ? AND sender_id != ?", roomID, false, readerID).
		Update("is_read", true).Error
}