package models

import (
	"time"
)


type ApprovedResearchCase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ResearchCaseID string 	 `gorm:"not null"`
	UniversityID   string  	 `gorm:"type:char(36);not null"`
	StudyProgramID string    `gorm:"type:char(36);not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	ResearchCase ResearchCase `gorm:"foreignKey:ResearchCaseID;"`
	University University `gorm:"foreignKey:UniversityID;references:ID;"`
	StudyProgram StudyProgram `gorm:"foreignKey:StudyProgramID;references:ID;"`
}


func (ApprovedResearchCase) TableName() string {
	return "ss_t_approved_research_case"
}
