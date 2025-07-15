package models

type StudentUser struct {
	UserID       string  `gorm:"primaryKey;type:char(36)"`
	UniversityID string  `gorm:"type:char(36);not null"`
	Jurusan      string  `gorm:"type:varchar(100);not null"`
	Nim          string  `gorm:"type:varchar(20);not null;unique"`
	Gpa          float64 `gorm:"type:decimal(3,2);not null"`

	User       User       `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	University University `gorm:"foreignKey:UniversityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (StudentUser) TableName() string {
	return "ss_student_user"
}
