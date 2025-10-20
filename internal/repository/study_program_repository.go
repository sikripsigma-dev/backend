package repository

import (
	"Skripsigma-BE/internal/models"
	"context"

	"gorm.io/gorm"
)

type StudyProgramRepository interface {
	FindAll(ctx context.Context) ([]models.StudyProgram, error)
	FindByID(ctx context.Context, id string) (*models.StudyProgram, error)
	Create(ctx context.Context, s *models.StudyProgram) error
	Update(ctx context.Context, s *models.StudyProgram) error
	Delete(ctx context.Context, id string) error
}

type studyProgramRepo struct {
	db *gorm.DB
}

func NewStudyProgramRepository(db *gorm.DB) StudyProgramRepository {
	return &studyProgramRepo{db}
}

func (r *studyProgramRepo) FindAll(ctx context.Context) ([]models.StudyProgram, error) {
	var items []models.StudyProgram
	err := r.db.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r *studyProgramRepo) FindByID(ctx context.Context, id string) (*models.StudyProgram, error) {
	var item models.StudyProgram
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *studyProgramRepo) Create(ctx context.Context, s *models.StudyProgram) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *studyProgramRepo) Update(ctx context.Context, s *models.StudyProgram) error {
	tx := r.db.WithContext(ctx).Model(&models.StudyProgram{}).
		Where("id = ?", s.ID).
		Updates(map[string]any{
			"name":       s.Name,
			"updated_at": s.UpdatedAt,
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *studyProgramRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&models.StudyProgram{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
