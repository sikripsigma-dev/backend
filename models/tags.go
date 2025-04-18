package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model Tag
type Tag struct {
	ID   string `gorm:"type:char(36);primaryKey" json:"id"`
	Name string `gorm:"type:varchar(255);unique;not null" json:"name"`
}

// Hook BeforeCreate untuk generate UUID sebelum insert ke database
func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}

// Nama tabel di database
func (Tag) TableName() string {
	return "ss_m_tags"
}
