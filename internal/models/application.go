package models

type Application struct {
	ID             uint   `gorm:"primaryKey"`
	ResearchCaseID string `gorm:"not null"`
	UserID         string `gorm:"not null"`
	Status         string `gorm:"default:'pending'"`
	AppliedAt      int64  `gorm:"autoCreateTime"`

	ResearchCase ResearchCase `gorm:"foreignKey:ResearchCaseID;constraint:OnDelete:CASCADE;"`
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (Application) TableName() string {
	return "ss_t_applications"
}