package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id       string `gorm:"type:char(36);primaryKey"`
	Nim      string `gorm:"type:varchar(50)"`
	Name     string
	Phone    string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	RoleId   uint   `gorm:"not null"`
	Status   string `gorm:"type:enum('active','inactive');default:'inactive'"`
	Image    string `gorm:"default:null"`
	Company  *CompanyUser `gorm:"foreignKey:UserID;references:Id"`
	Student  *StudentUser `gorm:"foreignKey:UserID;references:Id"`
	Supervisor *SupervisorUser `gorm:"foreignKey:UserID;references:Id"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Only generate new UUID if one isn't already set
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	return nil
}

func (User) TableName() string {
	return "ss_users"
}