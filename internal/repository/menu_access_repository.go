package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuAccessRepository interface {
	GetByRole(roleID uint) ([]models.MenuAccess, error)
	UpsertMany(rows []models.MenuAccess) error
	DeleteByRoleAndMenus(roleID uint, menuIDs []uint) error
}

type menuAccessRepository struct{ db *gorm.DB }

func NewMenuAccessRepository(db *gorm.DB) MenuAccessRepository {
	return &menuAccessRepository{db}
}

func (r *menuAccessRepository) GetByRole(roleID uint) ([]models.MenuAccess, error) {
	var res []models.MenuAccess
	err := r.db.Where("id_role = ?", roleID).Find(&res).Error
	return res, err
}

func (r *menuAccessRepository) UpsertMany(rows []models.MenuAccess) error {
	if len(rows) == 0 { return nil }
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id_role"}, {Name: "id_menu"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"can_view","can_add","can_edit","can_delete","can_approve","can_export",
		}),
	}).Create(&rows).Error
}

func (r *menuAccessRepository) DeleteByRoleAndMenus(roleID uint, menuIDs []uint) error {
	if len(menuIDs) == 0 { return nil }
	return r.db.Where("id_role = ? AND id_menu IN ?", roleID, menuIDs).
		Delete(&models.MenuAccess{}).Error
}
