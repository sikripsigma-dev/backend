package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type MonitoringProgressRepository interface {
	Create(feedback *models.SupervisorMonitoringProgress) error
	CompanyCreate(feedback *models.CompanyMonitoringProgress) error
	UpdateReportStatus(reportID uint, status string) error
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

