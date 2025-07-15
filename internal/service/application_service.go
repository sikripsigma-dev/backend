package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ApplicationService struct {
	applicationRepo         repository.ApplicationRepository
	roleRepo                repository.RoleRepository
	researchCaseRepository  repository.ResearchCaseRepository
	userRepository          repository.UserRepository
	assignmentRepository    repository.AssignmentRepository
}

func NewApplicationService(
	applicationRepo repository.ApplicationRepository,
	roleRepo repository.RoleRepository,
	researchCaseRepo repository.ResearchCaseRepository,
	userRepo repository.UserRepository,
	assignmentRepo repository.AssignmentRepository,
) *ApplicationService {
	return &ApplicationService{
		applicationRepo,
		roleRepo,
		researchCaseRepo,
		userRepo,
		assignmentRepo,
	}
}

func (s *ApplicationService) CreateApplication(req dto.CreateApplicationRequest, userID string) (*models.Application, error) {
	_, err := s.researchCaseRepository.GetByID(req.ResearchCaseID)
	if err != nil {
		return nil, fmt.Errorf("Invalid research_case_id")
	}

	application := &models.Application{
		ResearchCaseID: req.ResearchCaseID,
		UserID:         userID,
		Status:         "diajukan", // default status
	}

	if err := s.applicationRepo.Create(application); err != nil {
		return nil, err
	}

	return application, nil
}

func (s *ApplicationService) CheckApplicationExists(userID, researchCaseID string) (bool, error) {
	application, err := s.applicationRepo.GetByUserIDAndResearchCaseID(userID, researchCaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return application != nil, nil
}

func (s *ApplicationService) ProcessApplication(applicationID string, status string, userID string) error {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return fmt.Errorf("Application not found")
	}

	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	if user.Company == nil {
		return fmt.Errorf("You are not authorized to process this application")
	}

	application.Status = status
	application.ProcessedAt = time.Now()
	application.ProcessedBy = user.Company.UserID

	return s.applicationRepo.Update(application)
}

func (s *ApplicationService) RespondToApplication(applicationID string, status string, userID string) error {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return fmt.Errorf("Application not found")
	}

	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	if user.Student == nil {
		return fmt.Errorf("You are not authorized to process this application")
	}

	if application.UserID != userID {
		return fmt.Errorf("You are not allowed to respond to this application")
	}

	if status != "confirmed" && status != "declined" {
		return fmt.Errorf("Invalid status for student response")
	}

	application.Status = status
	application.ProcessedAt = time.Now()
	application.ProcessedBy = user.Student.UserID

	if err := s.applicationRepo.Update(application); err != nil {
		return err
	}

	if status == "confirmed" {
		// Cek apakah sudah ada assignment aktif
		assignments, err := s.assignmentRepository.GetByUserID(userID)
		if err != nil {
			return fmt.Errorf("Failed to check existing assignments: %v", err)
		}
		for _, a := range assignments {
			if a.Status == "active" {
				return fmt.Errorf("You already have an active assignment")
			}
		}

		// Cancel aplikasi lain
		err = s.applicationRepo.CancelOtherApplications(application.ID, userID)
		if err != nil {
			return fmt.Errorf("Failed to cancel other applications: %v", err)
		}

		// Buat assignment baru
		assignment := &models.Assignment{
			ApplicationID:  application.ID,
			UserID:         userID,
			ResearchCaseID: application.ResearchCaseID,
			Status:         "active",
			StartedAt:      time.Now(),
		}

		if err := s.assignmentRepository.Create(assignment); err != nil {
			return fmt.Errorf("Failed to create assignment: %v", err)
		}
	}

	return nil
}

func (s *ApplicationService) GetApplicationsByResearchCaseID(researchCaseID string) ([]models.Application, error) {
	applications, err := s.applicationRepo.GetByResearchCaseID(researchCaseID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get applications for research case: %v", err)
	}
	return applications, nil
}

func (s *ApplicationService) GetApplicationsByStudentID(studentID string) ([]models.Application, error) {
	applications, err := s.applicationRepo.GetAllByStudentID(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %v", err)
	}
	return applications, nil
}
