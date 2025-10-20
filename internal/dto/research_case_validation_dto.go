package dto

type CreateResearchCaseValidationRequest struct {
	ResearchCaseID string `json:"research_case_id" validate:"required"`
	ProgramHeadID  string `json:"program_head_id" validate:"required"`
	Status         string `json:"status" validate:"required,oneof=approved rejected pending"`
}

type CreateResearchCaseValidationCommentRequest struct {
	ResearchCaseID string `json:"research_case_id" validate:"required"`
	ProgramHeadID  string `json:"program_head_id" validate:"required"`
	Comment        string `json:"comment" validate:"required"`
}

// type UpdateResearchCaseValidationRequest struct {
// 	Status string `json:"status" validate:"required,oneof=approved rejected pending"`
// }

type UpdateResearchCaseValidationRequest struct {
	Status         string `json:"status" validate:"required,oneof=approved rejected pending"`
	ResearchCaseID string `json:"research_case_id,omitempty"` // ⬅️ baru (opsional)
}

type UpdateResearchCaseValidationCommentRequest struct {
	Comment string `json:"comment" validate:"required"`
}
