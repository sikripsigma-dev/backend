package models

// models/company_user.go
type CompanyUser struct {
	UserID      string `gorm:"primaryKey"`
	Division    string
	CompanyName string
	User        User   `gorm:"foreignKey:UserID;references:Id"`
}


func (CompanyUser) TableName() string {
	return "ss_company_user"
}