package models

import (
	"time"
)

type Assignment struct {
	ID             uint      `gorm:"primaryKey"`
	ApplicationID  uint      `gorm:"not null;unique"` // 1:1 dengan Application
	UserID         string    `gorm:"not null"`         // Untuk tracking siapa yang assigned
	ResearchCaseID string    `gorm:"not null"`         // Duplikat dari Application, agar query lebih cepat
	Status         string    `gorm:"default:'active'"` // active, completed, etc.
	StartedAt      time.Time `gorm:"autoCreateTime"`
	EndedAt        *time.Time

	Application  Application  `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE"`
	User         User         `gorm:"foreignKey:UserID"`
	ResearchCase ResearchCase `gorm:"foreignKey:ResearchCaseID"`
}

func (Assignment) TableName() string {
	return "ss_t_assignments"
}
