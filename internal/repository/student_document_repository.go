package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type StudentDocumentRepository interface {
	FindByUserID(userID string) ([]models.StudentDocument, error)
	Save(doc *models.StudentDocument) error
	Update(doc *models.StudentDocument) error
}

type studentDocumentRepository struct {
	db *gorm.DB
}

func NewStudentDocumentRepository(db *gorm.DB) StudentDocumentRepository {
	return &studentDocumentRepository{db}
}

func (r *studentDocumentRepository) FindByUserID(userID string) ([]models.StudentDocument, error) {
	var docs []models.StudentDocument
	err := r.db.Where("user_id = ?", userID).Preload("User").Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *studentDocumentRepository) Save(doc *models.StudentDocument) error {
	return r.db.Create(doc).Error
}

func (r *studentDocumentRepository) Update(doc *models.StudentDocument) error {
	return r.db.Save(doc).Error
}
