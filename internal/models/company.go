package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model Company
type Company struct {
	Id          string `gorm:"type:char(36);primaryKey"`
	Name        string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Phone       string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Description string
}

// Hook BeforeCreate untuk generate UUID sebelum insert ke database
func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	company.Id = uuid.New().String()
	return
}

// Tentukan nama tabel
func (Company) TableName() string {
	return "ss_m_companies"
}
