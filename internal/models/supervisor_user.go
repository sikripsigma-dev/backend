package models

// Dosen Pembimbing
type SupervisorUser struct {
	UserID       string `gorm:"primaryKey;type:char(36)"`
	UniversityID string `gorm:"type:char(36);not null"`
	Nidn         string `gorm:"type:varchar(20);not null;unique"`

	User       User       `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	University University `gorm:"foreignKey:UniversityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (SupervisorUser) TableName() string {
	return "ss_supervisor_user"
}