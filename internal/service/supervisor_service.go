package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/repository"
)

type SupervisorService struct {
	supervisorRepo repository.SupervisorRepository
}

func NewSupervisorService(supervisorRepo repository.SupervisorRepository) *SupervisorService {
	return &SupervisorService{supervisorRepo}
}

func (s *SupervisorService) GetStudentsBySupervisor(supervisorID string) ([]dto.StudentResponse, error) {
	return s.supervisorRepo.GetStudentsBySupervisor(supervisorID)
}
