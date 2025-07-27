package models

import (
	"time"
)

type SupervisorMonitoringProgress struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	WeeklyReportID  uint      `json:"weekly_report_id"` // FK ke laporan mingguan mahasiswa
	Feedback        string    `json:"feedback"`
	Status          string    `json:"status"`           // misal: "Diterima", "Revisi", dst.
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	WeeklyReport WeeklyReport `gorm:"foreignKey:WeeklyReportID;references:ID"` // relasi opsional
}

func (SupervisorMonitoringProgress) TableName() string {
	return "ss_t_supervisor_monitoring_progress"
}
