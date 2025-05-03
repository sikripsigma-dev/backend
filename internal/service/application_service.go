package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
)

type ApplicationService struct {
	applicationRepo         repository.ApplicationRepository
	roleRepo                repository.RoleRepository
	researchCaseRepository  repository.ResearchCaseRepository
}

func NewApplicationService(applicationRepo repository.ApplicationRepository, roleRepo repository.RoleRepository, researchCaseRepo repository.ResearchCaseRepository) *ApplicationService {
	return &ApplicationService{applicationRepo, roleRepo, researchCaseRepo}
}

func (s *ApplicationService) CreateApplication(req dto.CreateApplicationRequest) (*models.Application, error) {
	// Validasi apakah research case ada
	_, err := s.researchCaseRepository.GetByID(req.ResearchCaseID)
	if err != nil {
		return nil, fmt.Errorf("Invalid research_case_id")
	}

	// Konversi research_case_id dari string ke uint
	// researchCaseID, err := strconv.ParseUint(req.ResearchCaseID, 10, 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("Invalid research_case_id format")
	// }

	application := &models.Application{
		// ResearchCaseID: uint(researchCaseID),
		ResearchCaseID: req.ResearchCaseID,
		UserID:         req.UserID,
		Status:         req.Status,
	}

	if err := s.applicationRepo.Create(application); err != nil {
		return nil, err
	}

	return application, nil
}
