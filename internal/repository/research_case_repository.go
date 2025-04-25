package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type ResearchCaseRepository interface {
	Create(researchCase *models.ResearchCase) error
	GetByID(id string) (*models.ResearchCase, error)
	GetAll() ([]models.ResearchCase, error)
	AssociateTags(researchCaseID string, tagIDs []string) error
	FindByIDWithRelations(id string, out *models.ResearchCase) error
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

func (r *researchCaseRepository) GetByID(id string) (*models.ResearchCase, error) {
	var rc models.ResearchCase
	if err := r.db.First(&rc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rc, nil
}

func (r *researchCaseRepository) GetAll() ([]models.ResearchCase, error) {
	var cases []models.ResearchCase
	if err := r.db.Find(&cases).Error; err != nil {
		return nil, err
	}
	return cases, nil
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
		First(out, "id = ?", id).Error
}

