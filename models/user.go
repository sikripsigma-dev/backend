package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model User
type User struct {
	Id       string `gorm:"type:char(36);primaryKey"`
	Nim      string `gorm:"unique"`
	Name     string
	Phone	 string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

// Hook BeforeCreate untuk generate UUID sebelum insert ke database
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Id = uuid.New().String()
	return
}

// Tentukan nama tabel yang digunakan di database
func (User) TableName() string {
	return "ss_users"
}
