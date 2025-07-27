package models

// Kaprodi
type HeadstudyUser struct{
	UserID       string `gorm:"primaryKey;type:char(36)"`
	UniversityID string `gorm:"type:char(36);not null"`
	Nidn         string `gorm:"type:varchar(20);not null;unique"`
	StudyProgramID string `gorm:"type:char(36);not null"`

	User       User       `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	University University `gorm:"foreignKey:UniversityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StudyProgram StudyProgram `gorm:"foreignKey:StudyProgramID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (HeadstudyUser) TableName() string {
	return "ss_headstudy_user"
}