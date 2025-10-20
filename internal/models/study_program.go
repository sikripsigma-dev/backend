package models

import (
	"time"

	"gorm.io/gorm"
)

type StudyProgram struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null;uniqueIndex:ux_study_program_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (StudyProgram) TableName() string {
	 return "ss_m_study_program" 
}
