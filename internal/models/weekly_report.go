// models/report.go
package models

import (
	// "Skripsigma-BE/internal/util"
	"time"
)

type WeeklyReport struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID string    `json:"student_id"`
	Week      int       `json:"week"`
	Progress  string    `json:"progress"`
	Plans     string    `json:"plans"`
	Mood      int       `json:"mood"`
	Notes     string    `json:"notes"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Files     string    `json:"files"` // ← JSON string, bukan array langsung

	Student User `gorm:"foreignKey:StudentID;references:Id"`
}
func (WeeklyReport) TableName() string {
	return "ss_t_weekly_reports"
}