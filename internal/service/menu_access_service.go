package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type MenuAccessService interface {
	GetTreeByRole(roleID uint) ([]dto.MenuAccessNode, error)
	UpdateBulk(roleID uint, items []dto.UpdateMenuAccessItem) error
}

type menuAccessService struct {
	db         *gorm.DB
	menuRepo   repository.MenuRepository
	accessRepo repository.MenuAccessRepository
}

func NewMenuAccessService(db *gorm.DB, menuRepo repository.MenuRepository, accessRepo repository.MenuAccessRepository) MenuAccessService {
	return &menuAccessService{db: db, menuRepo: menuRepo, accessRepo: accessRepo}
}

func (s *menuAccessService) GetTreeByRole(roleID uint) ([]dto.MenuAccessNode, error) {
	menus, err := s.menuRepo.GetAll()
	if err != nil { return nil, err }

	accessRows, err := s.accessRepo.GetByRole(roleID)
	if err != nil { return nil, err }
	acc := map[uint]models.MenuAccess{}
	for _, a := range accessRows { acc[a.MenuID] = a }

	var build func(parent *uint) []dto.MenuAccessNode
	build = func(parent *uint) []dto.MenuAccessNode {
		var out []dto.MenuAccessNode
		for _, m := range menus {
			if (m.ParentID == nil && parent == nil) || (m.ParentID != nil && parent != nil && *m.ParentID == *parent) {
				a := acc[m.ID] // zero value jika tak ada baris
				node := dto.MenuAccessNode{
					ID:        m.ID,
					Name:      m.Nama,
					URL:       m.URL,
					Icon:      m.Icon,
					IsActive:  m.IsActive,
					ParentID:  m.ParentID,

					CanView:    a.CanView,
					CanAdd:     a.CanAdd,
					CanEdit:    a.CanEdit,
					CanDelete:  a.CanDelete,
					CanApprove: a.CanApprove,
					CanExport:  a.CanExport,
				}
				children := build(&m.ID)
				if len(children) > 0 { node.Children = children }
				out = append(out, node)
			}
		}
		return out
	}
	return build(nil), nil
}

func normalizeFlags(it *dto.UpdateMenuAccessItem) {
	if it.CanAdd || it.CanEdit || it.CanDelete || it.CanApprove || it.CanExport {
		it.CanView = true
	}
}

func (s *menuAccessService) UpdateBulk(roleID uint, items []dto.UpdateMenuAccessItem) error {
	if roleID == 0 { return errors.New("invalid role id") }

	var upserts []models.MenuAccess
	var deletes []uint

	for _, it := range items {
		if it.MenuID == 0 { continue }
		normalizeFlags(&it)

		if it.CanView || it.CanAdd || it.CanEdit || it.CanDelete || it.CanApprove || it.CanExport {
			upserts = append(upserts, models.MenuAccess{
				RoleID: roleID, MenuID: it.MenuID,
				CanView: it.CanView, CanAdd: it.CanAdd, CanEdit: it.CanEdit,
				CanDelete: it.CanDelete, CanApprove: it.CanApprove, CanExport: it.CanExport,
			})
		} else {
			deletes = append(deletes, it.MenuID)
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		ar := repository.NewMenuAccessRepository(tx)

		if len(upserts) > 0 {
			if err := ar.UpsertMany(upserts); err != nil { return err }
		}
		if len(deletes) > 0 {
			if err := ar.DeleteByRoleAndMenus(roleID, deletes); err != nil { return err }
		}
		return nil
	})
}
