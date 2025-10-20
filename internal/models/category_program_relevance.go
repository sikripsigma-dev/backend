// models/category_program_relevance.go
package models

import "time"

type CategoryProgramRelevance struct {
	ID             string    `gorm:"type:char(36);primaryKey"`
	CategoryID     string    `gorm:"type:char(36);not null;index:ux_cat_prog,unique"`
	StudyProgramID string    `gorm:"type:char(36);not null;index:ux_cat_prog,unique"`

	// Relevansi:
	// Weight 0–100 (semakin besar semakin relevan)
	Weight  uint8 `gorm:"type:tinyint unsigned;default:70" json:"weight"`
	// Penanda kategori inti/utama untuk prodi tsb
	IsCore  bool  `gorm:"default:false" json:"is_core"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Category     Category     `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StudyProgram StudyProgram `gorm:"foreignKey:StudyProgramID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CategoryProgramRelevance) TableName() string {
	return "ss_m_category_program_map"
}
