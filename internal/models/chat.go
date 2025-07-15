package models

import (
	"time"
)

type ChatRooms struct {
	ID        string    `gorm:"type:char(36);primaryKey"`
	StudentID string    `gorm:"type:char(36);not null"`
	CompanyID string    `gorm:"type:char(36);not null"` // tetap NOT NULL
	CreatedAt time.Time
	UpdatedAt time.Time

	UnreadCount int `gorm:"-"`

	Student User    `gorm:"foreignKey:StudentID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Company Company `gorm:"foreignKey:CompanyID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // GANTI DARI SET NULL KE CASCADE

	Messages []ChatMessage `gorm:"foreignKey:RoomID"`
}


type ChatMessage struct {
	ID        string    `gorm:"type:char(36);primaryKey"`
	RoomID    string    `gorm:"type:char(36);not null"`
	SenderID  *string   `gorm:"type:char(36)"` // bisa null
	SenderType string   `gorm:"type:enum('student','company');not null"`
	Message   string    `gorm:"type:text;not null"`
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relasi (jika mahasiswa)
	Sender  *User      `gorm:"foreignKey:SenderID;references:Id"`
	Room    ChatRooms  `gorm:"foreignKey:RoomID;references:ID"`
}


// Table names
func (ChatRooms) TableName() string {
	return "ss_t_chat_rooms"
}

func (ChatMessage) TableName() string {
	return "ss_t_chat_message"
}
