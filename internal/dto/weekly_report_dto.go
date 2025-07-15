// dto/weekly_report_dto.go
package dto

type CreateWeeklyReportDTO struct {
	Week      int      `form:"week" validate:"required"`
	Progress  string   `form:"progress" validate:"required"`
	Plans     string   `form:"plans" validate:"required"`
	Mood      int      `form:"mood" validate:"required"`
	Notes     string   `form:"notes"`
	StartDate string   `form:"start_date" validate:"required"`
	EndDate   string   `form:"end_date" validate:"required"`
	Files     []string `form:"-"` // tidak di-parse langsung, ditangani manual
}



type WeeklyReportResponse struct {
	ID        uint     `json:"id"`
	Week      int      `json:"week"`
	Progress  string   `json:"progress"`
	Plans     string   `json:"plans"`
	Mood      int      `json:"mood"`
	Notes     string   `json:"notes"`
	Status    string   `json:"status"`
	StartDate string   `json:"start_date"` 
	EndDate   string   `json:"end_date"` 
	Files     []string `json:"files"`
}
