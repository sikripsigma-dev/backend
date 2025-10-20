package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	Create(assignment *models.Assignment) error
	GetByUserID(userID string) ([]models.Assignment, error)
	GetActiveByUserID(userID string) (*models.Assignment, error)
	GetByCompanyID(companyID string) ([]models.Assignment, error)
	UpdateStatus(id uint, status string) error

	GetByIDForUser(id uint, userID string) (*models.Assignment, error)
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

// func (r *assignmentRepository) GetByUserID(userID string) ([]models.Assignment, error) {
// 	var assignments []models.Assignment
// 	err := r.db.Where("user_id = ?", userID).Find(&assignments).Error
// 	return assignments, err
// }

func (r *assignmentRepository) GetByUserID(userID string) ([]models.Assignment, error) {
    var assignments []models.Assignment
    err := r.db.
        Where("user_id = ?", userID).
        Preload("ResearchCase").
        Preload("ResearchCase.Company").
        Order("started_at DESC").
        Find(&assignments).Error
    return assignments, err
}

// func (r *assignmentRepository) GetActiveByUserID(userID string) (*models.Assignment, error) {
//     var assignment models.Assignment
//     err := r.db.
//         Preload("ResearchCase").
//         // Where("user_id = ? AND (status = ? OR status = ?)", userID, "active", "completed").
//         Where("user_id = ? AND (status = ? OR status = ?)", userID, "active").
//         Order("started_at DESC").
//         First(&assignment).Error
//     if err != nil {
//         return nil, err
//     }
//     return &assignment, nil
// }

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

func (r *assignmentRepository) GetByCompanyID(companyID string) ([]models.Assignment, error) {
	var assignments []models.Assignment
	err := r.db.
		Joins("JOIN ss_t_research_cases rc ON rc.id = ss_t_assignments.research_case_id").
		// Where("rc.company_id = ? AND ss_t_assignments.status = ?", companyID, "active").
		Where("rc.company_id = ?", companyID).
		Preload("User").
		Preload("ResearchCase").
		Preload("Application").
		Find(&assignments).Error

	return assignments, err
}

func (r *assignmentRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Assignment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *assignmentRepository) GetByIDForUser(id uint, userID string) (*models.Assignment, error) {
    var a models.Assignment
    err := r.db.
        Preload("ResearchCase").
        Preload("ResearchCase.Company").
        Where("id = ? AND user_id = ?", id, userID).
        First(&a).Error
    if err != nil {
        return nil, err
    }
    return &a, nil
}