package repository

import (
	"Skripsigma-BE/internal/models"
	"time"

	"gorm.io/gorm"
)

type ResearchCaseRepository interface {
	Create(researchCase *models.ResearchCase) error
	GetByID(id string) (*models.ResearchCase, error)
	GetAll() ([]models.ResearchCase, error)
	AssociateTags(researchCaseID string, tagIDs []string) error
	FindByIDWithRelations(id string, out *models.ResearchCase) error
	GetByCompanyID(companyID string) ([]models.ResearchCase, error)
	Update(researchCase *models.ResearchCase) error
	UpdateActiveStatus(id string, isActive bool) error
	GetAllForHeadstudy() ([]models.ResearchCase, error)
	GetAllForHeadstudyFiltered(programHeadID string) ([]models.ResearchCase, error)
	GetApprovedByHeadStudy(headStudyID string) ([]models.ResearchCase, error)
	AssociateCategories(researchCaseID string, categoryIDs []string) error
    ReplaceCategories(researchCaseID string, categoryIDs []string) error
}

type researchCaseRepository struct {
	db *gorm.DB
}

func NewResearchCaseRepository(db *gorm.DB) ResearchCaseRepository {
	return &researchCaseRepository{db}
}

func (r *researchCaseRepository) Create(researchCase *models.ResearchCase) error {
	return r.db.Create(researchCase).Error
}

// func (r *researchCaseRepository) GetByID(id string) (*models.ResearchCase, error) {
// 	var rc models.ResearchCase
// 	if err := r.db.Preload("Company").First(&rc, "id = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &rc, nil
// }

func (r *researchCaseRepository) GetByID(id string) (*models.ResearchCase, error) {
    var rc models.ResearchCase
    if err := r.db.
        Preload("Company").
        Preload("Tags").
        Preload("Categories"). // ⬅️ tambah
        First(&rc, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &rc, nil
}

// func (r *researchCaseRepository) GetAll() ([]models.ResearchCase, error) {
// 	var researchCases []models.ResearchCase
// 	err := r.db.
// 		Where("is_active = ?", true).
// 		Where("activated_by_admin = ?", true).
// 		Preload("Company").
// 		Preload("Tags").
// 		Find(&researchCases).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return researchCases, nil
// }

func (r *researchCaseRepository) GetAll() ([]models.ResearchCase, error) {
    var researchCases []models.ResearchCase
    err := r.db.
        Where("is_active = ?", true).
        Where("activated_by_admin = ?", true).
        Preload("Company").
        Preload("Tags").
        Preload("Categories"). // ⬅️ tambah
        Find(&researchCases).Error
    return researchCases, err
}

func (r *researchCaseRepository) GetApprovedByHeadStudy(studentID string) ([]models.ResearchCase, error) {
	var student models.StudentUser
	// ambil data student beserta HeadStudyID
	if err := r.db.First(&student, "user_id = ?", studentID).Error; err != nil {
		return nil, err
	}

	var researchCases []models.ResearchCase
	err := r.db.
		Model(&models.ResearchCase{}).
		Joins("JOIN ss_t_research_case_validations rcv ON rcv.research_case_id = ss_t_research_cases.id").
		Where("ss_t_research_cases.is_active = ?", true).
		// Where("ss_t_research_cases.activated_by_admin = ?", true).
		Where("rcv.program_head_id = ?", student.HeadStudyID).
		Where("rcv.status = ?", "approved").
		Preload("Company").
		Preload("Tags").
		Find(&researchCases).Error


	if err != nil {
		return nil, err
	}
	return researchCases, nil
}



// func (r *researchCaseRepository) GetAllForHeadstudy() ([]models.ResearchCase, error) {
// 	var researchCases []models.ResearchCase
// 	err := r.db.
// 		Where("is_active = ?", true).
// 		Where("activated_by_admin = ?", false).
// 		Preload("Company").
// 		Preload("Tags").
// 		Find(&researchCases).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return researchCases, nil
// }

func (r *researchCaseRepository) GetAllForHeadstudy() ([]models.ResearchCase, error) {
    var researchCases []models.ResearchCase
    err := r.db.
        Where("is_active = ?", true).
        Where("activated_by_admin = ?", false).
        Preload("Company").
        Preload("Tags").
        Preload("Categories"). // ⬅️ tambah
        Find(&researchCases).Error
    return researchCases, err
}

func (r *researchCaseRepository) AssociateTags(researchCaseID string, tagIDs []string) error {
	for _, tagID := range tagIDs {
		researchCaseTag := models.ResearchCaseTag{
			ResearchCaseID: researchCaseID,
			TagID:          tagID,
		}
		if err := r.db.Create(&researchCaseTag).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *researchCaseRepository) FindByIDWithRelations(id string, out *models.ResearchCase) error {
    return r.db.
        Preload("Company").
        Preload("Tags").
        Preload("Categories").
        First(out, "id = ?", id).Error
}

// get by company id
func (r *researchCaseRepository) GetByCompanyID(companyID string) ([]models.ResearchCase, error) {
	var cases []models.ResearchCase
	if err := r.db.Where("company_id = ?", companyID).Find(&cases).Error; err != nil {
		return nil, err
	}
	return cases, nil
}

// update research case
func (r *researchCaseRepository) Update(researchCase *models.ResearchCase) error {
	return r.db.Save(researchCase).Error
}

func (r *researchCaseRepository) UpdateActiveStatus(id string, isActive bool) error {
	return r.db.Model(&models.ResearchCase{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (r *researchCaseRepository) AssociateCategories(researchCaseID string, categoryIDs []string) error {
    if len(categoryIDs) == 0 { return nil }
    rows := make([]models.ResearchCaseCategory, 0, len(categoryIDs))
    now := time.Now()
    for _, cid := range categoryIDs {
        rows = append(rows, models.ResearchCaseCategory{
            ResearchCaseID: researchCaseID,
            CategoryID:     cid,
            CreatedAt:      now,
        })
    }
    return r.db.Create(&rows).Error
}

func (r *researchCaseRepository) ReplaceCategories(researchCaseID string, categoryIDs []string) error {
    if err := r.db.
        Where("research_case_id = ?", researchCaseID).
        Delete(&models.ResearchCaseCategory{}).Error; err != nil {
        return err
    }
    return r.AssociateCategories(researchCaseID, categoryIDs)
}

func (r *researchCaseRepository) GetAllForHeadstudyFiltered(programHeadID string) ([]models.ResearchCase, error) {
    var researchCases []models.ResearchCase

    // Hanya yang aktif & belum dipublish admin
    // Join:
    //  - rcc: m2m research_case ⟷ category
    //  - cpr: mapping category ⟷ study_program
    //  - hs : ambil StudyProgramID kaprodi dari user_id
    err := r.db.
        Table("ss_t_research_cases").
        Joins("JOIN ss_t_research_case_categories rcc ON rcc.research_case_id = ss_t_research_cases.id").
        Joins("JOIN ss_m_category_program_map cpr ON cpr.category_id = rcc.category_id").
        Joins("JOIN ss_headstudy_user hs ON hs.user_id = ? AND hs.study_program_id = cpr.study_program_id", programHeadID).
        Where("ss_t_research_cases.is_active = ? AND ss_t_research_cases.activated_by_admin = ?", true, false).
        Group("ss_t_research_cases.id").
        Preload("Company").
        Preload("Tags").
        Preload("Categories").
        Find(&researchCases).Error

    return researchCases, err
}
