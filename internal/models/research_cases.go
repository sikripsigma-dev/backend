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
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relasi ke Company
	Company Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`

	// Relasi ke Tags (many-to-many)
	Tags []Tag `gorm:"many2many:ss_t_research_case_tags;foreignKey:ID;joinForeignKey:ResearchCaseID;References:ID;joinReferences:TagID" json:"tags"`
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
