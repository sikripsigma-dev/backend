package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type WeeklyReportRepository interface {
	Create(report *models.WeeklyReport) error
	GetByStudent(studentID string) ([]models.WeeklyReport, error)
	GetByUniversity(universityID string) ([]models.WeeklyReport, error)
}

type weeklyReportRepo struct {
	db *gorm.DB
}

func NewWeeklyReportRepository(db *gorm.DB) WeeklyReportRepository {
	return &weeklyReportRepo{db}
}

func (r *weeklyReportRepo) Create(report *models.WeeklyReport) error {
	return r.db.Create(report).Error
}

func (r *weeklyReportRepo) GetByStudent(studentID string) ([]models.WeeklyReport, error) {
	var reports []models.WeeklyReport
	err := r.db.
		Where("student_id = ?", studentID).
		Order("week desc").
		Preload("Student").
		Find(&reports).Error
	return reports, err
}

func (r *weeklyReportRepo) GetByUniversity(universityID string) ([]models.WeeklyReport, error) {
	var reports []models.WeeklyReport
	err := r.db.
		Joins("JOIN ss_student_user su ON su.user_id = ss_t_weekly_reports.student_id").
		Joins("JOIN ss_users u ON u.id = su.user_id").
		Where("su.university_id = ?", universityID).
		Preload("Student").
		Find(&reports).Error

	return reports, err
}
