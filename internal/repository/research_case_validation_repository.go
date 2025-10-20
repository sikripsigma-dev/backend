// repository/research_case_validation_repository.go
package repository

import (
	"errors"
	"time"

	"Skripsigma-BE/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResearchCaseValidationRepository interface {
    // validations
    CreateValidation(v *models.ResearchCaseValidation) error
    GetValidationsByResearchCaseID(researchCaseID string) ([]models.ResearchCaseValidation, error)
    GetValidationsByResearchCaseAndProgramHeadID(researchCaseID, programHeadID string) ([]models.ResearchCaseValidation, error)

    // ⬇️ baru
    FindOne(researchCaseID, programHeadID string) (*models.ResearchCaseValidation, error)
    UpdateStatus(id, status string) error

    // comments
    CreateComment(c *models.ResearchCaseValidationComment) error
    GetCommentsByResearchCaseID(researchCaseID string) ([]models.ResearchCaseValidationComment, error)

    // ⬇️ approved helpers
    EnsureApproved(researchCaseID, programHeadID string) error
    RemoveApproved(researchCaseID, programHeadID string) error
}

type researchCaseValidationRepository struct{ db *gorm.DB }
func NewResearchCaseValidationRepository(db *gorm.DB) ResearchCaseValidationRepository {
    return &researchCaseValidationRepository{db}
}

func (r *researchCaseValidationRepository) CreateValidation(v *models.ResearchCaseValidation) error {
    v.ID = uuid.New().String()
    return r.db.Create(v).Error
}

func (r *researchCaseValidationRepository) GetValidationsByResearchCaseAndProgramHeadID(researchCaseID, programHeadID string) ([]models.ResearchCaseValidation, error) {
    var out []models.ResearchCaseValidation
    err := r.db.Preload("ProgramHead").
        Where("research_case_id = ? AND program_head_id = ?", researchCaseID, programHeadID).
        Find(&out).Error
    return out, err
}

func (r *researchCaseValidationRepository) GetValidationsByResearchCaseID(researchCaseID string) ([]models.ResearchCaseValidation, error) {
    var out []models.ResearchCaseValidation
    err := r.db.Preload("ProgramHead").
        Where("research_case_id = ?", researchCaseID).
        Find(&out).Error
    return out, err
}

// ⬇️ baru
func (r *researchCaseValidationRepository) FindOne(researchCaseID, programHeadID string) (*models.ResearchCaseValidation, error) {
    var v models.ResearchCaseValidation
    err := r.db.Where("research_case_id = ? AND program_head_id = ?", researchCaseID, programHeadID).
        First(&v).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, err
        }
        return nil, err
    }
    return &v, nil
}

// ⬇️ baru
func (r *researchCaseValidationRepository) UpdateStatus(id, status string) error {
    return r.db.Model(&models.ResearchCaseValidation{}).
        Where("id = ?", id).
        Updates(map[string]any{"status": status, "validated_at": time.Now()}).Error
}

// func (r *researchCaseValidationRepository) CreateComment(c *models.ResearchCaseValidationComment) error {
//     c.ID = uuid.New().String()
//     return r.db.Create(c).Error
// }

func (r *researchCaseValidationRepository) CreateComment(c *models.ResearchCaseValidationComment) error {
    c.ID = uuid.New().String()
    if err := r.db.Create(c).Error; err != nil {
        return err
    }
    // preload supaya response langsung berisi ProgramHead
    return r.db.Preload("ProgramHead").First(c, "id = ?", c.ID).Error
}


func (r *researchCaseValidationRepository) GetCommentsByResearchCaseID(researchCaseID string) ([]models.ResearchCaseValidationComment, error) {
    var out []models.ResearchCaseValidationComment
    err := r.db.Preload("ProgramHead").
        Where("research_case_id = ?", researchCaseID).
        Order("created_at DESC").
        Find(&out).Error
    return out, err
}

// ===== ApprovedResearchCase helpers =====

func (r *researchCaseValidationRepository) EnsureApproved(researchCaseID, programHeadID string) error {
    // ambil Univ & Prodi dari Kaprodi
    var hs models.HeadstudyUser
    if err := r.db.Where("user_id = ?", programHeadID).First(&hs).Error; err != nil {
        return err
    }
    ar := models.ApprovedResearchCase{
        ResearchCaseID: researchCaseID,
        UniversityID:   hs.UniversityID,
        StudyProgramID: hs.StudyProgramID,
    }
    // upsert sederhana
    return r.db.
        Where("research_case_id = ? AND university_id = ? AND study_program_id = ?",
            ar.ResearchCaseID, ar.UniversityID, ar.StudyProgramID).
        FirstOrCreate(&ar).Error
}

func (r *researchCaseValidationRepository) RemoveApproved(researchCaseID, programHeadID string) error {
    var hs models.HeadstudyUser
    if err := r.db.Where("user_id = ?", programHeadID).First(&hs).Error; err != nil {
        // kalau kaprodi record tak ada, dianggap non-fatal
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil }
        return err
    }
    return r.db.
        Where("research_case_id = ? AND university_id = ? AND study_program_id = ?",
            researchCaseID, hs.UniversityID, hs.StudyProgramID).
        Delete(&models.ApprovedResearchCase{}).Error
}
