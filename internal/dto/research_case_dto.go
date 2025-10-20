package dto

import (
	"Skripsigma-BE/internal/models"
	"errors"
	"strings"
)

type CreateResearchCaseRequest struct {
	CompanyID   string `json:"company_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Field       string `json:"field" validate:"required"`
	Location   string `json:"location" validate:"required"`
	Duration	string `json:"duration" validate:"required"`
	EducationRequirement string `json:"education_requirement" validate:"required"`
	TagIDs     []string `json:"tag_ids"`
	CategoryIDs          []string `json:"category_ids"`
}

type UpdateResearchCaseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Field       string `json:"field" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Duration    string `json:"duration" validate:"required"`
	EducationRequirement string `json:"education_requirement" validate:"required"`
	CategoryIDs          []string `json:"category_ids"`
}

type ResearchCaseWithValidationStatus struct {
	models.ResearchCase
	Reviewed     bool   `json:"reviewed"`
	ReviewStatus string `json:"review_status"` // "pending", "approved", "rejected", or "" if not reviewed
}


func (r *CreateResearchCaseRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(r.Description) == "" {
		return errors.New("description is required")
	}
	if strings.TrimSpace(r.Field) == "" {
		return errors.New("field is required")
	}
	return nil
}