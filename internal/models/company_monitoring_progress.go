package models

import (
	"time"
)

type CompanyMonitoringProgress struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	WeeklyReportID uint      `json:"weekly_report_id"` // FK ke laporan mingguan
	Feedback       string    `json:"feedback"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relasi ke laporan mingguan perusahaan
	WeeklyReport CompanyWeeklyReport `gorm:"foreignKey:WeeklyReportID;references:ID"`
}


func (CompanyMonitoringProgress) TableName() string {
	return "ss_t_company_monitoring_progress"
}
