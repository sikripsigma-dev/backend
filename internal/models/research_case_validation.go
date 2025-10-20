package models

import (
	"time"
)

type ResearchCaseValidation struct {
	ID             string    `gorm:"type:char(36);primaryKey"`
	ResearchCaseID string    `gorm:"type:char(36)"`
	ProgramHeadID  string    `gorm:"type:char(36)"`
	Status         string    `gorm:"type:enum('approved','rejected','pending');default:'pending'"`
	ValidatedAt    time.Time `gorm:""`

	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	ProgramHead   User         `gorm:"foreignKey:ProgramHeadID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ResearchCase  ResearchCase `gorm:"foreignKey:ResearchCaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}


type ResearchCaseValidationComment struct {
	ID             string    `gorm:"type:char(36);primaryKey"`
	ResearchCaseID string    `gorm:"type:char(36);"`
	ProgramHeadID  string    `gorm:"type:char(36);"`
	Comment        string    `gorm:"type:text;"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	ProgramHead   User         `gorm:"foreignKey:ProgramHeadID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ResearchCase  ResearchCase `gorm:"foreignKey:ResearchCaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}


func (ResearchCaseValidation) TableName() string {
	return "ss_t_research_case_validations"
}

func (ResearchCaseValidationComment) TableName() string {
	return "ss_t_research_case_validation_comments"
}