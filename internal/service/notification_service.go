package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type NotificationService interface {
	Create(notification *models.Notification) error
	GetForReceiver(userID, companyID string) ([]models.Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID, companyID string) error
	GetByID(id string) (*models.Notification, error)
	CreateNotification(notif *models.Notification) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo}
}

func (s *notificationService) Create(notification *models.Notification) error {
	return s.notificationRepo.Create(notification)
}

func (s *notificationService) GetForReceiver(userID, companyID string) ([]models.Notification, error) {
	return s.notificationRepo.GetForReceiver(userID, companyID)
}

func (s *notificationService) GetByID(id string) (*models.Notification, error) {
	return s.notificationRepo.GetByID(id)
}

func (s *notificationService) MarkAsRead(id string) error {
	return s.notificationRepo.MarkAsRead(id)
}

func (s *notificationService) MarkAllAsRead(userID, companyID string) error {
	return s.notificationRepo.MarkAllAsRead(userID, companyID)
}

func (s *notificationService) CreateNotification(notif *models.Notification) error {
	return s.notificationRepo.Create(notif)
}
