package models

import (
	"time"
)

type CompanyWeeklyReport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	StudentID      string    `json:"student_id"`
	ResearchCaseID string    `json:"research_case_id"`
	Week           int       `json:"week"`
	Activities     string    `json:"activities"`
	Issues         string    `json:"issues"`
	Hopes          string    `json:"hopes"`
	Notes          string    `json:"notes"`
	Files          string    `json:"files"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`

	Student      User         `gorm:"foreignKey:StudentID;references:Id"`
	ResearchCase ResearchCase `gorm:"foreignKey:ResearchCaseID;references:ID"`
}

func (CompanyWeeklyReport) TableName() string {
	return "ss_t_company_weekly_reports"
}
