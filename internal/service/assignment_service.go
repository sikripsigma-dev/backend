package service

import (
	// "Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/repository"
	"fmt"
)

type AssignmentService struct{
	assignmentRepo repository.AssignmentRepository
}

func NewAssignmentService(assignmentRepo repository.AssignmentRepository) * AssignmentService{
	return &AssignmentService{assignmentRepo}
}

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

func (s *AssignmentService) GetActiveAssignmentsByCompany(companyID string) ([]dto.AssignmentResponse, error) {
	assignments, err := s.assignmentRepo.GetByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	var responses []dto.AssignmentResponse
	for _, a := range assignments {
		resp := dto.AssignmentResponse{
			ID:             a.ID,
			ApplicationID:  a.ApplicationID,
			UserID:         a.UserID,
			ResearchCaseID: a.ResearchCaseID,
			Status:         a.Status,
			StartedAt:      a.StartedAt,
			EndedAt:        a.EndedAt,
		}

		// Tambahkan info research case
		if a.ResearchCase.ID != "" {
			resp.ResearchCase = &dto.AssignmentResearchCaseResponse{
				ID:                   a.ResearchCase.ID,
				CompanyID:            a.ResearchCase.CompanyID,
				Title:                a.ResearchCase.Title,
				Field:                a.ResearchCase.Field,
				Location:             a.ResearchCase.Location,
				EducationRequirement: a.ResearchCase.EducationRequirement,
				Duration:             a.ResearchCase.Duration,
				Description:          a.ResearchCase.Description,
				CreatedAt:            a.ResearchCase.CreatedAt,
			}
		}

		// Tambahkan info user
		if a.User.Id != "" {
			resp.User = &dto.UserResponse{
				Id:    a.User.Id,
				Name:  a.User.Name,
				Email: a.User.Email,
			}
		}

		// Tambahkan info company
		if a.ResearchCase.Company.Id != "" {
			resp.ResearchCase.Company = &dto.AssignmentCompanyResponse{
				ID:          a.ResearchCase.Company.Id,
				Name:        a.ResearchCase.Company.Name,
				Email:       a.ResearchCase.Company.Email,
				Phone:       a.ResearchCase.Company.Phone,
				Address:     a.ResearchCase.Company.Address,
				Description: a.ResearchCase.Company.Description,
			}
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *AssignmentService) UpdateAssignmentStatus(id uint, status string) error {
	if status != "active" && status != "inactive" {
		return fmt.Errorf("Invalid status: %s", status)
	}
	return s.assignmentRepo.UpdateStatus(id, status)
}
