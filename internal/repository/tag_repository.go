package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *models.Tag) error
	GetByID(id string) (*models.Tag, error)
	GetAll() ([]models.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) GetByID(id string) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.Where("id = ?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}