package models

import "time"

type UserCreateLog struct {
	ID         string    `gorm:"type:char(36);primaryKey"`
	CreatedBy  string    `gorm:"type:char(36);not null"`  // ID admin pembuat akun
	UserID     string    `gorm:"type:char(36);not null"`  // ID user yang dibuat
	RoleID     uint      `gorm:"not null"`
	Email      string    `gorm:"not null"`
	Name       string    `gorm:"not null"`
	CreatedAt  time.Time
}

func (UserCreateLog) TableName() string {
	return "ss_user_create_log"
}