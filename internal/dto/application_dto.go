package dto

import (
	// "Skripsigma-BE/internal/models"
	"time"
)

type CreateApplicationRequest struct {
	ResearchCaseID string `json:"research_case_id" validate:"required"`
	// UserID         string `json:"user_id" validate:"required"`
	// Status         string `json:"status" validate:"required"`
}

type ProcessApplicationRequest struct {
	Status string `json:"status" validate:"required"`
}

// type ApplicationResponse struct {
// 	ID             uint   `json:"id"`
// 	ResearchCaseID string `json:"research_case_id"`
// 	UserID         string `json:"user_id"`
// 	Status         string `json:"status"`
// 	AppliedAt      int64  `json:"applied_at"`
// }

type ApplicationResponse struct {
	ID             uint                    `json:"id"`
	ResearchCaseID string                  `json:"research_case_id"`
	Status         string                  `json:"status"`
	AppliedAt      time.Time                   `json:"applied_at"`
	ProcessedAt    time.Time                   `json:"processed_at"`
	ProcessedBy    string                  `json:"processed_by"`
	User           ApplicationUserResponse `json:"user"`
}

type ApplyResponse struct {
	ID             uint   `json:"id"`
	ResearchCaseID string `json:"research_case_id"`
	UserID         string `json:"user_id"`
	Status         string `json:"status"`
	AppliedAt      time.Time `json:"applied_at"`
	ProcessedAt    time.Time `json:"processed_at"`
}

type ApplicationUserResponse struct {
	Id    string `json:"id"`
	Nim   string `json:"nim"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type GetApplicationByStudentResponse struct {
	ID                  uint                 `json:"id"`
	Status              string               `json:"status"`
	AppliedAt           time.Time            `json:"applied_at"`
	ProcessedAt         *time.Time           `json:"processed_at,omitempty"`
	ProcessedBy         string               `json:"processed_by,omitempty"`
	ResearchCase        ResearchCaseResponse `json:"research_case"`
}

type ResearchCaseResponse struct {
	ID                   string          `json:"id"`
	Title                string          `json:"title"`
	Field                string          `json:"field"`
	Location             string          `json:"location"`
	EducationRequirement string          `json:"education_requirement"`
	Duration             string          `json:"duration"`
	Description          string          `json:"description"`
	CreatedAt            time.Time       `json:"created_at"`
	Tags                 []TagResponse   `json:"tags"`
	Company              CompanyResponse `json:"company"`
}

type CompanyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
