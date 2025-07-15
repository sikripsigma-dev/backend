package models

import (
	"time"

	"gorm.io/gorm"
)

type University struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (University) TableName() string {
	return "ss_m_university"
}
