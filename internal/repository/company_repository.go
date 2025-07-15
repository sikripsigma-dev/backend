package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	Create(company *models.Company) error
	GetByID(id string) (*models.Company, error)
	GetAll() ([]models.Company, error)
	Update(company *models.Company) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db}
}

func (r *companyRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) GetByID(id string) (*models.Company, error) {
	var company models.Company
	if err := r.db.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetAll() ([]models.Company, error) {
	var companies []models.Company
	if err := r.db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

// update company
func (r *companyRepository) Update(company *models.Company) error {
	return r.db.Save(company).Error
}

