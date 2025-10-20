package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResearchCase model
type ResearchCase struct {
	ID                   string    `gorm:"type:char(36);primaryKey" json:"id"`
	CompanyID            string    `gorm:"type:char(36);not null" json:"company_id"`
	Title                string    `gorm:"not null" json:"title"`
	Field                string    `json:"field"`
	Location 		     string    `json:"location"`
	EducationRequirement string    `json:"education_requirement"`
	Duration             string    `json:"duration"`
	Description          string    `json:"description"`
	IsActive 			 bool 	   `gorm:"default:true" json:"is_active"`
	ActivatedByAdmin   	 bool 	   `gorm:"default:true" json:"activated_by_admin"`
	ActivatedByAdminAt   time.Time `gorm:"default:null" json:"activated_by_admin_at"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	ResearchCaseValidations []ResearchCaseValidation `gorm:"foreignKey:ResearchCaseID;references:ID" json:"research_case_validations"`
	ResearchCaseValidationComments []ResearchCaseValidationComment `gorm:"foreignKey:ResearchCaseID;references:ID" json:"research_case_validation_comments"`

	// Relasi ke Company
	Company Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`

	// Relasi ke Tags (many-to-many)
	Tags []Tag `gorm:"many2many:ss_t_research_case_tags;foreignKey:ID;joinForeignKey:ResearchCaseID;References:ID;joinReferences:TagID" json:"tags"`

	Categories  []Category  `gorm:"many2many:ss_t_research_case_categories;foreignKey:ID;joinForeignKey:ResearchCaseID;References:ID;joinReferences:CategoryID" json:"categories"`
}

// models/research_case_category.go
type ResearchCaseCategory struct {
    ResearchCaseID string    `gorm:"type:char(36);primaryKey"`
    CategoryID     string    `gorm:"type:char(36);primaryKey"`
    CreatedAt      time.Time `gorm:"autoCreateTime"`
}

// Hook BeforeCreate untuk generate UUID sebelum insert ke database
func (r *ResearchCase) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}

// Nama tabel di database
func (ResearchCase) TableName() string {
	return "ss_t_research_cases"
}

func (ResearchCaseCategory) TableName() string { 
	return "ss_t_research_case_categories" 
}
