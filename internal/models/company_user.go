package models

type CompanyUser struct {
	UserID    string `gorm:"primaryKey;type:char(36)"`
	CompanyID string `gorm:"type:char(36);not null"`
	Division  string

	User    User    `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Company Company `gorm:"foreignKey:CompanyID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CompanyUser) TableName() string {
	return "ss_company_user"
}