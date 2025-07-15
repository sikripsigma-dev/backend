package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetByUserID(userID string) ([]models.Notification, error)
	GetForReceiver(userID, companyID string) ([]models.Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID, companyID string) error
	GetByID(id string) (*models.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

// Get by user
func (r *notificationRepository) GetByUserID(userID string) ([]models.Notification, error) {
	var notifs []models.Notification
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifs).Error
	return notifs, err
}

// Get by user or company (global handler)
func (r *notificationRepository) GetForReceiver(userID, companyID string) ([]models.Notification, error) {
	var notifs []models.Notification
	err := r.db.
		Where("user_id = ? OR company_id = ?", userID, companyID).
		Order("created_at DESC").
		Find(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) GetByID(id string) (*models.Notification, error) {
	var notif models.Notification
	err := r.db.First(&notif, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}

// Mark as read
func (r *notificationRepository) MarkAsRead(id string) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID, companyID string) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? OR company_id = ?", userID, companyID).
		Update("is_read", true).Error
}