package dto

type CreateSupervisorFeedbackDTO struct {
	WeeklyReportID uint   `json:"weekly_report_id" validate:"required"`
	Feedback       string `json:"feedback" validate:"required"`
	Status         string `json:"status" validate:"required"`
}
type CreateCompanyFeedbackDTO struct {
	WeeklyReportID uint   `json:"weekly_report_id" validate:"required"`
	Feedback       string `json:"feedback" validate:"required"`
}
