package repository

import (
	"Skripsigma-BE/internal/models"
	"context"

	"gorm.io/gorm"
)

type UniversityRepository interface {
	FindAll(ctx context.Context) ([]models.University, error)
	FindByID(ctx context.Context, id string) (*models.University, error)
	Create(ctx context.Context, u *models.University) error
	Update(ctx context.Context, u *models.University) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) UniversityRepository {
	return &repository{db}
}

func (r *repository) FindAll(ctx context.Context) ([]models.University, error) {
	var universities []models.University
	err := r.db.WithContext(ctx).Find(&universities).Error
	return universities, err
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.University, error) {
	var university models.University
	err := r.db.WithContext(ctx).First(&university, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &university, nil
}

func (r *repository) Create(ctx context.Context, u *models.University) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *repository) Update(ctx context.Context, u *models.University) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.University{}, "id = ?", id).Error
}
