package models

import (
	"time"
)

type Application struct {
	ID             uint   `gorm:"primaryKey"`
	ResearchCaseID string `gorm:"not null"`
	UserID         string `gorm:"not null"`
	Status         string `gorm:"default:'pending'"`
	AppliedAt      time.Time  `gorm:"autoCreateTime" json:"applied_at"`
	ProcessedAt time.Time `gorm:"autoUpdateTime" json:"processed_at"`
	ProcessedBy    string `gorm:"default:null"`

	ResearchCase ResearchCase `gorm:"foreignKey:ResearchCaseID;constraint:OnDelete:CASCADE;"`
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (Application) TableName() string {
	return "ss_t_applications"
}