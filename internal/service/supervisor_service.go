package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"

	"gorm.io/gorm"
)

type SupervisorService struct {
	supervisorRepo repository.SupervisorRepository
	DB             *gorm.DB
}

func NewSupervisorService(supervisorRepo repository.SupervisorRepository, db *gorm.DB) *SupervisorService {
	return &SupervisorService{supervisorRepo: supervisorRepo, DB: db}
}

func (s *SupervisorService) GetStudentsBySupervisor(supervisorID string) ([]dto.StudentResponse, error) {
	return s.supervisorRepo.GetStudentsBySupervisor(supervisorID)
}

func (s *SupervisorService) GetAllSupervisorsByHeadstudy(userID uint) ([]models.SupervisorUser, error) {
	var studyProgramUser models.HeadstudyUser
	if err := s.DB.Where("user_id = ?", userID).First(&studyProgramUser).Error; err != nil {
		return nil, err
	}

	var supervisors []models.SupervisorUser
	if err := s.DB.Preload("User").Where("university_id = ?", studyProgramUser.UniversityID).Find(&supervisors).Error; err != nil {
		return nil, err
	}

	return supervisors, nil
}
