package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/repository"
)

type HeadstudyService interface {
	GetStudentsByUniversity(univID string) ([]dto.StudentResponse, error)
}

type headstudyService struct {
	repo repository.HeadstudyRepository
}

func NewHeadstudyService(repo repository.HeadstudyRepository) HeadstudyService {
	return &headstudyService{
		repo: repo,
	}
}

func (s *headstudyService) GetStudentsByUniversity(univID string) ([]dto.StudentResponse, error) {
	return s.repo.GetStudentsByUniversity(univID)
}
