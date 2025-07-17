package models

type StudentDocument struct {
	ID     string `gorm:"type:char(36);primaryKey"`
	UserID string `gorm:"type:char(36)"`

	User User
	Type string `gorm:"type:text"`
	Path string `gorm:"type:text"`
}

func (StudentDocument) TableName() string {
	return "ss_student_document"
}
