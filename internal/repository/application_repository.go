package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(application *models.Application) error
	GetByID(id string) (*models.Application, error)
	Update(application *models.Application) error
	// GetAll() ([]models.Application, error)
	// GetByStudentID(studentID string) (*models.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db}
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(application).Error
}

// func (r *applicationRepository) GetAll() ([]models.Application, error) {
	
// 	var applications []models.Application
// 	if err := r.db.Find(&applications).Error; err != nil {
// 		return nil, err
// 	}
// 	return applications, nil
// }

func (r *applicationRepository) GetByID(id string) (*models.Application, error) {
	
	var application models.Application
	if err := r.db.Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) Update(application *models.Application) error {
	return r.db.Save(application).Error
}
