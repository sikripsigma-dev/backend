// service/research_case_validation_service.go
package service

import (
	"errors"
	"strings"
	"time"

	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"

	"gorm.io/gorm"
)

type ResearchCaseValidationService struct {
    researchCaseValidationRepo repository.ResearchCaseValidationRepository
}

func NewResearchCaseValidationService(repo repository.ResearchCaseValidationRepository) *ResearchCaseValidationService {
    return &ResearchCaseValidationService{researchCaseValidationRepo: repo}
}

func (s *ResearchCaseValidationService) CreateValidation(req dto.CreateResearchCaseValidationRequest) (*models.ResearchCaseValidation, error) {
    v := &models.ResearchCaseValidation{
        ResearchCaseID: req.ResearchCaseID,
        ProgramHeadID:  req.ProgramHeadID,
        Status:         normalize(req.Status),
    }
    if err := s.researchCaseValidationRepo.CreateValidation(v); err != nil {
        return nil, err
    }
    return v, nil
}

func (s *ResearchCaseValidationService) GetValidationsByResearchCaseID(id string) ([]models.ResearchCaseValidation, error) {
    return s.researchCaseValidationRepo.GetValidationsByResearchCaseID(id)
}

func (s *ResearchCaseValidationService) CreateComment(req dto.CreateResearchCaseValidationCommentRequest) (*models.ResearchCaseValidationComment, error) {
    c := &models.ResearchCaseValidationComment{
        ResearchCaseID: req.ResearchCaseID,
        ProgramHeadID:  req.ProgramHeadID,
        Comment:        strings.TrimSpace(req.Comment),
    }
    if err := s.researchCaseValidationRepo.CreateComment(c); err != nil {
        return nil, err
    }
    return c, nil
}

func (s *ResearchCaseValidationService) GetCommentsByResearchCaseID(id string) ([]models.ResearchCaseValidationComment, error) {
    return s.researchCaseValidationRepo.GetCommentsByResearchCaseID(id)
}

// ⬇️ baru: upsert + sinkron ApprovedResearchCase
func (s *ResearchCaseValidationService) UpsertValidationStatus(programHeadID, researchCaseID, status string) error {
    status = normalize(status)
    if !allowed(status) { return errors.New("invalid status") }

    existing, err := s.researchCaseValidationRepo.FindOne(researchCaseID, programHeadID)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { return err }

    if existing != nil && existing.ID != "" {
        if err := s.researchCaseValidationRepo.UpdateStatus(existing.ID, status); err != nil {
            return err
        }
    } else {
        v := &models.ResearchCaseValidation{
            ResearchCaseID: researchCaseID,
            ProgramHeadID:  programHeadID,
            Status:         status,
            ValidatedAt:    time.Now(),
        }
        if err := s.researchCaseValidationRepo.CreateValidation(v); err != nil {
            return err
        }
    }

    // sinkron tabel approved
    if status == "approved" {
        return s.researchCaseValidationRepo.EnsureApproved(researchCaseID, programHeadID)
    }
    // pending/rejected → hapus mapping approved
    return s.researchCaseValidationRepo.RemoveApproved(researchCaseID, programHeadID)
}

func normalize(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
func allowed(s string) bool {
    return s == "approved" || s == "rejected" || s == "pending"
}
