package dto

type CreateApplicationRequest struct {
	ResearchCaseID string `json:"research_case_id" validate:"required"`
	UserID         string `json:"user_id" validate:"required"`
	Status         string `json:"status" validate:"required"`
}

type ApplicationResponse struct {
	ID             uint   `json:"id"`
	ResearchCaseID string `json:"research_case_id"`
	UserID         string `json:"user_id"`
	Status         string `json:"status"`
	AppliedAt      int64  `json:"applied_at"`
}