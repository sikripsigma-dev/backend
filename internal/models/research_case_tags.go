package models

// ResearchCaseTag (Pivot Table)
type ResearchCaseTag struct {
	ResearchCaseID string      `gorm:"type:char(36);not null;primaryKey"`
	TagID          string      `gorm:"type:char(36);not null;primaryKey"`
	ResearchCase   ResearchCase `gorm:"foreignKey:ResearchCaseID;constraint:OnDelete:CASCADE"`
	Tag            Tag          `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// Nama tabel di database
func (ResearchCaseTag) TableName() string {
	return "ss_t_research_case_tags"
}
