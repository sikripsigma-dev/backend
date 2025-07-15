package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type ChatRepository interface {
	GetChatRoomsByStudentID(studentID string) ([]models.ChatRooms, error)
	GetMessagesByRoomID(roomID string) ([]models.ChatMessage, error)
	GetChatRoomsByCompanyID(companyID string) ([]models.ChatRooms, error)
}

// Struct implementasi
type chatRepository struct {
	db *gorm.DB
}

// Constructor
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db}
}

// Implementasi fungsi interface
func (r *chatRepository) GetChatRoomsByStudentID(studentID string) ([]models.ChatRooms, error) {
	var chatRooms []models.ChatRooms
	err := r.db.
		Preload("Company").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(1)
		}).
		Where("student_id = ?", studentID).
		Find(&chatRooms).Error

	if err != nil {
		return nil, err
	}

	for i := range chatRooms {
		var count int64
		r.db.Model(&models.ChatMessage{}).
			Where("room_id = ? AND is_read = ? AND sender_type != ?", chatRooms[i].ID, false, "student").
			Count(&count)
		chatRooms[i].UnreadCount = int(count)
	}

	return chatRooms, err
}

func (r *chatRepository) GetChatRoomsByCompanyID(companyID string) ([]models.ChatRooms, error) {
	var chatRooms []models.ChatRooms
	err := r.db.
		Preload("Company").
		Preload("Messages", func(db *gorm.DB) *gorm.DB{
			return db.Order("created_at desc").Limit(1)
		}).
		Where("company_id = ?", companyID).
		Find(&chatRooms).Error

	if err != nil {
		return nil, err
	}

	for i := range chatRooms {
		var count int64
		r.db.Model(&models.ChatMessage{}).
			Where("room_id = ? AND is_read = ? AND sender_type != ?", chatRooms[i].ID, false, "company").
			Count(&count)
		chatRooms[i].UnreadCount = int(count)
	}

	return chatRooms, err
}

func (r *chatRepository) GetMessagesByRoomID(roomID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage

	err := r.db.
		Where("room_id = ?", roomID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func SaveChatMessage(db *gorm.DB, msg *models.ChatMessage) error {
	return db.Create(msg).Error
}

func MarkMessagesAsRead(db *gorm.DB, roomID string, receiverType string) error {
	return db.Model(&models.ChatMessage{}).
		Where("room_id = ? AND sender_type != ? AND is_read = ?", roomID, receiverType, false).
		Update("is_read", true).Error
}
