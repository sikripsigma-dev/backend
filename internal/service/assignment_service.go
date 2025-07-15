package service

import (
	// "Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/repository"
)

type AssignmentService struct{
	assignmentRepo repository.AssignmentRepository
}

func NewAssignmentService(assignmentRepo repository.AssignmentRepository) * AssignmentService{
	return &AssignmentService{assignmentRepo}
}

// func (s *AssignmentService) GetActiveAssignment(userID string) (*models.Assignment, error) {
// 	return s.assignmentRepo.GetActiveByUserID(userID)
// }

func (s *AssignmentService) GetActiveAssignment(userID string) (*dto.AssignmentResponse, error) {
	assignment, err := s.assignmentRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := &dto.AssignmentResponse{
		ID:             assignment.ID,
		ApplicationID:  assignment.ApplicationID,
		UserID:         assignment.UserID,
		ResearchCaseID: assignment.ResearchCaseID,
		Status:         assignment.Status,
		StartedAt:      assignment.StartedAt,
		EndedAt:        assignment.EndedAt,
	}

	if assignment.ResearchCase.ID != "" {
		resp.ResearchCase = &dto.AssignmentResearchCaseResponse{
			ID:                   assignment.ResearchCase.ID,
			CompanyID:            assignment.ResearchCase.CompanyID,
			Title:                assignment.ResearchCase.Title,
			Field:                assignment.ResearchCase.Field,
			Location:             assignment.ResearchCase.Location,
			EducationRequirement: assignment.ResearchCase.EducationRequirement,
			Duration:             assignment.ResearchCase.Duration,
			Description:          assignment.ResearchCase.Description,
			CreatedAt:            assignment.ResearchCase.CreatedAt,
		}

		// Tambahkan jika company ada
		if assignment.ResearchCase.Company.Id != "" {
			resp.ResearchCase.Company = &dto.AssignmentCompanyResponse{
				ID:          assignment.ResearchCase.Company.Id,
				Name:        assignment.ResearchCase.Company.Name,
				Email:       assignment.ResearchCase.Company.Email,
				Phone:       assignment.ResearchCase.Company.Phone,
				Address:     assignment.ResearchCase.Company.Address,
				Description: assignment.ResearchCase.Company.Description,
			}
		}
	}

	return resp, nil
}

