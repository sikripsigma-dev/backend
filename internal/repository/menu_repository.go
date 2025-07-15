package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(menu *models.Menu) error
	GetByID(id string) (*models.Menu, error)
	Update(menu *models.Menu) error
	GetAll() ([]models.Menu, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

func (r * menuRepository) Update(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) GetAll() ([]models.Menu, error) {
	var menus []models.Menu
	if err := r.db.Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) GetByID(id string) (*models.Menu, error) {
	var menu models.Menu
	if err := r.db.Where("id = ?", id).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}