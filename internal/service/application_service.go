package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
	"time"
)

type ApplicationService struct {
	applicationRepo         repository.ApplicationRepository
	roleRepo                repository.RoleRepository
	researchCaseRepository  repository.ResearchCaseRepository
	userRepository           repository.UserRepository
}

func NewApplicationService(applicationRepo repository.ApplicationRepository, 
	roleRepo repository.RoleRepository, 
	researchCaseRepo repository.ResearchCaseRepository, 
	userRepo repository.UserRepository,
	) *ApplicationService {
	return &ApplicationService{
		applicationRepo, 
		roleRepo, 
		researchCaseRepo, 
		userRepo} 
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

func (s *ApplicationService) ProcessApplication(applicationID string, status string, userID string) error {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return fmt.Errorf("Application not found")
	}

	user, err := s.userRepository.GetWithCompanyByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	if user.Company == nil {
		return fmt.Errorf("You are not authorized to process this application")
	}

	application.Status = status
	application.ProcessedAt = time.Now().Unix()
	application.ProcessedBy = user.Company.UserID // atau user.Company.UserID tergantung desain

	return s.applicationRepo.Update(application)
}
