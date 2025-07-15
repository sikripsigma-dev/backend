package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
    ID        uint           `gorm:"primaryKey;column:id_menu"`
    Nama      string         `gorm:"column:nama;type:varchar(100);not null" json:"name"` 
    URL       *string        `gorm:"column:url;type:varchar(255)" json:"url,omitempty"`
    Icon      string         `gorm:"column:icon;type:varchar(255)" json:"icon"`
    IsActive  bool           `gorm:"column:is_active;default:true" json:"is_active"`
    ParentID  *uint          `gorm:"column:parent_id"`
    CreatedAt *time.Time     `gorm:"column:created_at"`
    UpdatedAt *time.Time     `gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Menu) TableName() string {
    return "ss_m_menu"
}
