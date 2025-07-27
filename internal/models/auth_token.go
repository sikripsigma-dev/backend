package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthToken struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	Type      string    `json:"type"` // e.g. "email_verification"
	IsUsed    bool      `json:"is_used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the table name used by GORM
func (AuthToken) TableName() string {
	return "ss_auth_tokens"
}

// BeforeCreate auto-generates UUID for ID field before insertion
func (t *AuthToken) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
