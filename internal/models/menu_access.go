package models

type MenuAccess struct {
	ID     uint `gorm:"primaryKey"`
	RoleID uint `gorm:"column:id_role;not null;index"`
	MenuID uint `gorm:"column:id_menu;not null;index"`

	CanView    bool `gorm:"default:false"`
	CanAdd     bool `gorm:"default:false"`
	CanEdit    bool `gorm:"default:false"`
	CanDelete  bool `gorm:"default:false"`
	CanApprove bool `gorm:"default:false"`
	CanExport  bool `gorm:"default:false"`
}

func (MenuAccess) TableName() string {
	return "ss_m_akses_menu"
}
