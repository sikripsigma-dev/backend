package dto

import (
	"time"
)

type AssignmentCompanyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type AssignmentResearchCaseResponse struct {
	ID                   string                    `json:"id"`
	CompanyID            string                    `json:"company_id"`
	Title                string                    `json:"title"`
	Field                string                    `json:"field"`
	Location             string                    `json:"location"`
	EducationRequirement string                    `json:"education_requirement"`
	Duration             string                    `json:"duration"`
	Description          string                    `json:"description"`
	CreatedAt            time.Time                 `json:"created_at"`
	Company              *AssignmentCompanyResponse `json:"company,omitempty"`
}

type AssignmentResponse struct {
	ID             uint                        `json:"id"`
	ApplicationID  uint                        `json:"application_id"`
	UserID         string                      `json:"user_id"`
	ResearchCaseID string                      `json:"research_case_id"`
	Status         string                      `json:"status"`
	StartedAt      time.Time                   `json:"started_at"`
	EndedAt        *time.Time                  `json:"ended_at,omitempty"`
	ResearchCase   *AssignmentResearchCaseResponse `json:"research_case,omitempty"`
}
