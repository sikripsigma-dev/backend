package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	Id          string `gorm:"type:char(36);primaryKey"`
	Name        string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Phone       string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Description string
	Industry    string `gorm:"not null"`
	Website     string
	Status      string `gorm:"type:enum('active','inactive');default:'active'"`
	Logo        string
}

func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	company.Id = uuid.New().String()
	return
}

func (Company) TableName() string {
	return "ss_m_companies"
}
