// internal/models/category.go
package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(120);not null;uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// (opsional) siapkan relasi ke research case via M:N
	// Join table akan kita buat terpisah: ss_t_research_case_categories
	ResearchCases []ResearchCase `gorm:"many2many:ss_t_research_case_categories;foreignKey:ID;joinForeignKey:CategoryID;References:ID;joinReferences:ResearchCaseID" json:"-"`
}

func (Category) TableName() string { return "ss_m_categories" }

// --- Hooks (ID & slug) ---

func (c *Category) BeforeCreate(_ any) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if strings.TrimSpace(c.Slug) == "" && c.Name != "" {
		c.Slug = Slugify(c.Name)
	}
	return nil
}

func (c *Category) BeforeSave(_ any) (err error) {
	// kalau slug kosong tapi name ada → generate
	if strings.TrimSpace(c.Slug) == "" && c.Name != "" {
		c.Slug = Slugify(c.Name)
	}
	return nil
}

// util sederhana; boleh pindah ke util package
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// ganti spasi & underscore dengan strip
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// rapikan strip ganda
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	// (opsional) buang karakter non-alnum/- jika mau lebih ketat
	return s
}
