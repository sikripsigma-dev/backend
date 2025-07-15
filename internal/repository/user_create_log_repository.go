package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type UserCreateLogRepository interface {
	Create(log *models.UserCreateLog) error
	GetDB() *gorm.DB
}

type userCreateLogRepository struct {
	db *gorm.DB
}

func NewUserCreateLogRepository(db *gorm.DB) UserCreateLogRepository {
	return &userCreateLogRepository{db}
}

func (r *userCreateLogRepository) Create(log *models.UserCreateLog) error {
	return r.db.Create(log).Error
}

func (r *userCreateLogRepository) GetDB() *gorm.DB {
    return r.db
}