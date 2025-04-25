package dto

import (
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