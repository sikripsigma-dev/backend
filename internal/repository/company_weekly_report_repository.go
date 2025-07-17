package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type CompanyWeeklyReportRepository interface {
	Create(report *models.CompanyWeeklyReport) error
	GetByStudent(studentID, researchCaseID string) ([]models.CompanyWeeklyReport, error)
}

type companyWeeklyReportRepo struct {
	db *gorm.DB
}

func NewCompanyWeeklyReportRepository(db *gorm.DB) CompanyWeeklyReportRepository {
	return &companyWeeklyReportRepo{db}
}

func (r *companyWeeklyReportRepo) Create(report *models.CompanyWeeklyReport) error {
	return r.db.Create(report).Error
}

func (r *companyWeeklyReportRepo) GetByStudent(studentID string, researchCaseID string) ([]models.CompanyWeeklyReport, error) {
	var reports []models.CompanyWeeklyReport
	err := r.db.
		Where("student_id = ?", studentID).
		Where("research_case_id = ?", researchCaseID).
		Order("week desc").
		Preload("Student").
		Preload("ResearchCase").
		Find(&reports).Error
	return reports, err
}
