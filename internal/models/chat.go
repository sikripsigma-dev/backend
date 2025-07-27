package models

import (
	"time"
)

type ChatRooms struct {
	ID              string    `gorm:"type:char(36);primaryKey"`
	InitiatorID     string    `gorm:"type:char(36);not null"`
	TargetUserID    *string   `gorm:"type:char(36)"`
	TargetCompanyID *string   `gorm:"type:char(36)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Initiator     User     `gorm:"foreignKey:InitiatorID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TargetUser    *User    `gorm:"foreignKey:TargetUserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TargetCompany *Company `gorm:"foreignKey:TargetCompanyID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Messages    []ChatMessage `gorm:"foreignKey:RoomID"`
	UnreadCount int           `gorm:"-"`
}

func (ChatRooms) TableName() string {
	return "ss_t_chat_rooms"
}

type ChatMessage struct {
	ID         string    `gorm:"type:char(36);primaryKey"`
	RoomID     string    `gorm:"type:char(36);not null"`
	SenderID   *string   `gorm:"type:char(36)"`
	SenderType string    `gorm:"type:enum('student','company','supervisor');not null"`
	Message    string    `gorm:"type:text;not null"`
	IsRead     bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Sender *User     `gorm:"foreignKey:SenderID;references:Id"`
	Room   ChatRooms `gorm:"foreignKey:RoomID;references:ID"`
}

func (ChatMessage) TableName() string {
	return "ss_t_chat_message"
}