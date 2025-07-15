package models

import (
	"time"

	"Skripsigma-BE/internal/util"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Notification struct {
	ID         string         `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     *string        `gorm:"type:char(36)" json:"user_id,omitempty"`
	CompanyID  *string        `gorm:"type:char(36)" json:"company_id,omitempty"`
	Type       string         `gorm:"type:varchar(50);not null" json:"type"`
	Message    string         `gorm:"type:text;not null" json:"message"`
	Metadata   util.JSONB     `gorm:"type:json" json:"metadata"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}


// BeforeCreate hook untuk generate UUID sebelum insert ke database
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.New().String()
	return nil
}

func (Notification) TableName() string {
	return "ss_t_notification"
}
