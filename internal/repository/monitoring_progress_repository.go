package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type MonitoringProgressRepository interface {
	Create(feedback *models.SupervisorMonitoringProgress) error
	CompanyCreate(feedback *models.CompanyMonitoringProgress) error
	UpdateReportStatus(reportID uint, status string) error
	GetSupervisorFeedbackByReportID(reportID uint) ([]models.SupervisorMonitoringProgress, error)
	GetCompanyFeedbackByReportID(reportID uint) ([]models.CompanyMonitoringProgress, error)
}

type monitoringProgressRepository struct {
	db *gorm.DB
}

func NewMonitoringProgressRepository(db *gorm.DB) MonitoringProgressRepository {
	return &monitoringProgressRepository{db}
}

func (r *monitoringProgressRepository) Create(feedback *models.SupervisorMonitoringProgress) error {
	return r.db.Create(feedback).Error
}

func (r *monitoringProgressRepository) CompanyCreate(feedback *models.CompanyMonitoringProgress) error {
	return r.db.Create(feedback).Error
}

func (r *monitoringProgressRepository) UpdateReportStatus(reportID uint, status string) error {
	return r.db.Model(&models.WeeklyReport{}).
		Where("id = ?", reportID).
		Update("status", status).Error
}

func (r *monitoringProgressRepository) GetSupervisorFeedbackByReportID(reportID uint) ([]models.SupervisorMonitoringProgress, error) {
	var feedbacks []models.SupervisorMonitoringProgress
	err := r.db.
		Where("weekly_report_id = ?", reportID).
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *monitoringProgressRepository) GetCompanyFeedbackByReportID(reportID uint) ([]models.CompanyMonitoringProgress, error) {
	var feedbacks []models.CompanyMonitoringProgress
	err := r.db.
		Where("weekly_report_id = ?", reportID).
		Find(&feedbacks).Error
	return feedbacks, err
}

