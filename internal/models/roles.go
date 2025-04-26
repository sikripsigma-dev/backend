package models

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

func (Role) TableName() string {
	return "ss_m_roles"
}