package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	Create(assignment *models.Assignment) error
	GetByUserID(userID string) ([]models.Assignment, error)
	GetActiveByUserID(userID string) (*models.Assignment, error)
}

type assignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return &assignmentRepository{db}
}

func (r *assignmentRepository) Create(assignment *models.Assignment) error {
	return r.db.Create(assignment).Error
}

func (r *assignmentRepository) GetByUserID(userID string) ([]models.Assignment, error) {
	var assignments []models.Assignment
	err := r.db.Where("user_id = ?", userID).Find(&assignments).Error
	return assignments, err
}

func (r *assignmentRepository) GetActiveByUserID(userID string) (*models.Assignment, error) {
	var assignment models.Assignment
	err := r.db.
		Preload("ResearchCase").
		Where("user_id = ? AND status = ?", userID, "active").
		First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}