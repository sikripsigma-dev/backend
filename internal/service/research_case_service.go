package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
)

type ResearchCaseService struct {
	researchCaseRepo             repository.ResearchCaseRepository
	researchCaseValidationRepo  repository.ResearchCaseValidationRepository
}

func NewResearchCaseService(repo repository.ResearchCaseRepository, validationRepo repository.ResearchCaseValidationRepository) *ResearchCaseService {
	return &ResearchCaseService{
		researchCaseRepo:            repo,
		researchCaseValidationRepo:  validationRepo,
	}
}

type ResearchCaseWithValidationStatus struct {
	models.ResearchCase
	Reviewed     bool   `json:"reviewed"`
	ReviewStatus string `json:"review_status"`
}

func (s *ResearchCaseService) CreateResearchCase(req dto.CreateResearchCaseRequest) (*models.ResearchCase, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid request: %v", err)
	}

	researchCase := models.ResearchCase{
		CompanyID:             req.CompanyID,
		Title:                 req.Title,
		Description:           req.Description,
		Field:                 req.Field,
		Location:              req.Location,
		Duration:              req.Duration,
		EducationRequirement:  req.EducationRequirement,
	}

	if err := s.researchCaseRepo.Create(&researchCase); err != nil {
		return nil, fmt.Errorf("Failed to create research case: %v", err)
	}

	if len(req.TagIDs) > 0 {
		if err := s.researchCaseRepo.AssociateTags(researchCase.ID, req.TagIDs); err != nil {
			return nil, fmt.Errorf("Failed to associate tags: %v", err)
		}
	}

	if len(req.CategoryIDs) > 0 {
		if err := s.researchCaseRepo.AssociateCategories(researchCase.ID, req.CategoryIDs); err != nil {
			return nil, fmt.Errorf("failed to associate categories: %v", err)
		}
	}

	var fullResearchCase models.ResearchCase
	if err := s.researchCaseRepo.FindByIDWithRelations(researchCase.ID, &fullResearchCase); err != nil {
		return nil, fmt.Errorf("Failed to fetch created research case with relations: %v", err)
	}

	return &fullResearchCase, nil
}

func (s *ResearchCaseService) GetAllResearchCases() ([]models.ResearchCase, error) {
	researchCases, err := s.researchCaseRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch research cases: %v", err)
	}
	return researchCases, nil
}

func (s *ResearchCaseService) GetApprovedByHeadStudy(studentID string) ([]models.ResearchCase, error) {
	researchCases, err := s.researchCaseRepo.GetApprovedByHeadStudy(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch research cases for student: %v", err)
	}
	return researchCases, nil
}


// func (s *ResearchCaseService) GetAllResearchCasesForHeadstudy(programHeadID string) ([]ResearchCaseWithValidationStatus, error) {
// 	researchCases, err := s.researchCaseRepo.GetAllForHeadstudy()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch research cases: %v", err)
// 	}

// 	var results []ResearchCaseWithValidationStatus

// 	for _, rc := range researchCases {
// 		validations, err := s.researchCaseValidationRepo.GetValidationsByResearchCaseAndProgramHeadID(rc.ID, programHeadID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to fetch validation for research_case_id %s: %v", rc.ID, err)
// 		}

// 		status := ""
// 		reviewed := false
// 		if len(validations) > 0 {
// 			status = validations[0].Status
// 			reviewed = true
// 		}

// 		results = append(results, ResearchCaseWithValidationStatus{
// 			ResearchCase: rc,
// 			Reviewed:     reviewed,
// 			ReviewStatus: status,
// 		})
// 	}

// 	return results, nil
// }

func (s *ResearchCaseService) GetResearchCaseByID(id string) (*models.ResearchCase, error) {
	researchCase, err := s.researchCaseRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Research case not found: %v", err)
	}
	return researchCase, nil
}

func (s *ResearchCaseService) GetResearchCasesByCompanyID(companyID string) ([]models.ResearchCase, error) {
	researchCases, err := s.researchCaseRepo.GetByCompanyID(companyID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch research cases for company: %v", err)
	}
	return researchCases, nil
}

func (s *ResearchCaseService) UpdateResearchCase(id string, req dto.UpdateResearchCaseRequest) (*models.ResearchCase, error) {
	researchCase, err := s.researchCaseRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Research case not found: %v", err)
	}

	researchCase.Title = req.Title
	researchCase.Description = req.Description
	researchCase.Field = req.Field
	researchCase.Location = req.Location
	researchCase.Duration = req.Duration
	researchCase.EducationRequirement = req.EducationRequirement

	if err := s.researchCaseRepo.Update(researchCase); err != nil {
		return nil, fmt.Errorf("Failed to update research case: %v", err)
	}

	if err := s.researchCaseRepo.ReplaceCategories(researchCase.ID, req.CategoryIDs); err != nil {
		return nil, fmt.Errorf("failed to replace categories: %v", err)
	}

	return researchCase, nil
}

func (s *ResearchCaseService) SetResearchCaseActiveStatus(id string, isActive bool) error {
	_, err := s.researchCaseRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("Research case not found: %v", err)
	}
	return s.researchCaseRepo.UpdateActiveStatus(id, isActive)
}

func (s *ResearchCaseService) GetAllResearchCasesForHeadstudy(programHeadID string) ([]ResearchCaseWithValidationStatus, error) {
    // ⬇️ sebelumnya: researchCases, err := s.researchCaseRepo.GetAllForHeadstudy()
    researchCases, err := s.researchCaseRepo.GetAllForHeadstudyFiltered(programHeadID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch research cases: %v", err)
    }

    var results []ResearchCaseWithValidationStatus
    for _, rc := range researchCases {
        vals, err := s.researchCaseValidationRepo.GetValidationsByResearchCaseAndProgramHeadID(rc.ID, programHeadID)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch validation for research_case_id %s: %v", rc.ID, err)
        }
        reviewed := len(vals) > 0
        status := ""
        if reviewed { status = vals[0].Status }

        results = append(results, ResearchCaseWithValidationStatus{
            ResearchCase: rc,
            Reviewed:     reviewed,
            ReviewStatus: status,
        })
    }
    return results, nil
}

