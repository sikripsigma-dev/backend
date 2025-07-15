package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type UserCreateLogService interface {
	CreateLog(log *models.UserCreateLog) error
	GetLogs() ([]models.UserCreateLog, error)
    GetLogsByCreator(creatorID string) ([]models.UserCreateLog, error)
    GetLogsByUser(userID string) ([]models.UserCreateLog, error)       
}

type userCreateLogService struct {
	repo repository.UserCreateLogRepository
}

func NewUserCreateLogService(repo repository.UserCreateLogRepository) UserCreateLogService {
	return &userCreateLogService{repo}
}

func (s *userCreateLogService) CreateLog(log *models.UserCreateLog) error {
	return s.repo.Create(log)
}

func (s *userCreateLogService) GetLogs() ([]models.UserCreateLog, error) {
    var logs []models.UserCreateLog
    err := s.repo.GetDB().Find(&logs).Error
    return logs, err
}

func (s *userCreateLogService) GetLogsByCreator(creatorID string) ([]models.UserCreateLog, error) {
    var logs []models.UserCreateLog
    err := s.repo.GetDB().Where("created_by = ?", creatorID).Find(&logs).Error
    return logs, err
}

func (s *userCreateLogService) GetLogsByUser(userID string) ([]models.UserCreateLog, error) {
    var logs []models.UserCreateLog
    err := s.repo.GetDB().Where("user_id = ?", userID).Find(&logs).Error
    return logs, err
}
