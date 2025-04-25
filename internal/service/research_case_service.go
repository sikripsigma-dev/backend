package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
)

type ResearchCaseService struct {
	researchCaseRepo repository.ResearchCaseRepository
}

func NewResearchCaseService(researchCaseRepo repository.ResearchCaseRepository) *ResearchCaseService {
	return &ResearchCaseService{researchCaseRepo}
}

func (s *ResearchCaseService) CreateResearchCase(req dto.CreateResearchCaseRequest) (*models.ResearchCase, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid request: %v", err)
	}

		researchCase := models.ResearchCase{
			CompanyID: 				req.CompanyID,
			Title:       			req.Title,
			Description: 			req.Description,
			Field:      			req.Field,
			Location:    			req.Location,
			Duration:    			req.Duration,
			EducationRequirement: 	req.EducationRequirement,
		}

	if err := s.researchCaseRepo.Create(&researchCase); err != nil {
		return nil, fmt.Errorf("Failed to create research case: %v", err)
	}

	if len(req.TagIDs) > 0 {
		if err := s.researchCaseRepo.AssociateTags(researchCase.ID, req.TagIDs); err != nil {
			return nil, fmt.Errorf("Failed to associate tags: %v", err)
		}
	}

	var fullResearchCase models.ResearchCase
	if err := s.researchCaseRepo.FindByIDWithRelations(researchCase.ID, &fullResearchCase); err != nil {
		return nil, fmt.Errorf("Failed to fetch created research case with relations: %v", err)
	}

	// return &researchCase, nil
	return &fullResearchCase, nil
}

func (s *ResearchCaseService) GetAllResearchCases() ([]models.ResearchCase, error) {
	researchCases, err := s.researchCaseRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch research cases: %v", err)
	}
	return researchCases, nil
}

func (s *ResearchCaseService) GetResearchCaseByID(id string) (*models.ResearchCase, error) {
	researchCase, err := s.researchCaseRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Research case not found: %v", err)
	}
	return researchCase, nil
}


