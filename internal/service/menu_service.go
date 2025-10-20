package service

import (
	"Skripsigma-BE/internal/constants"
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type MenuService struct {
	menuRepository     repository.MenuRepository
	assignmentRepo     repository.AssignmentRepository
	DB                 *gorm.DB
}

func NewMenuService(
	menuRepository repository.MenuRepository,
	assignmentRepo repository.AssignmentRepository,
	db *gorm.DB,
) *MenuService {
	return &MenuService{
		menuRepository: menuRepository,
		assignmentRepo: assignmentRepo,
		DB:             db,
	}
}

// func (s *MenuService) CreateMenu(req dto.CreateMenuRequest) (*models.Menu, error) {
// 	menu := &models.Menu{
// 		Nama:    	req.Name,
// 		URL:         &req.URL,
// 		ParentID: req.ParentID,
// 	}

// 	if err := s.menuRepository.Create(menu); err != nil {
// 		return nil, err
// 	}
// 	return menu, nil

// }

// func (s *MenuService) UpdateMenu(ID string, req dto.UpdateMenuRequest) (*models.Menu, error) {
// 	menu, err := s.menuRepository.GetByID(ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("Menu Not Found")
// 	}

// 	menu.Nama = req.Name
// 	menu.URL = &req.URL
// 	menu.IsActive = req.IsActive
// 	menu.ParentID = req.ParentID

// 	if err := s.menuRepository.Update(menu); err != nil {
// 		return nil, err
// 	}

// 	return menu, nil
// }


func (s *MenuService) CreateMenu(req dto.CreateMenuRequest) (*models.Menu, error) {
	// (Opsional) validasi parent ada
	if req.ParentID != nil {
		if _, err := s.menuRepository.GetByID(*req.ParentID); err != nil {
			return nil, fmt.Errorf("parent menu tidak ditemukan")
		}
	}
	menu := &models.Menu{
		Nama:     req.Name,
		URL:      &req.URL,   // req.URL sudah divalidasi via dto tag
		Icon:     req.Icon,
		IsActive: req.IsActive, // default di model true, tapi pakai input kalau ada
		ParentID: req.ParentID,
	}
	if err := s.menuRepository.Create(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *MenuService) UpdateMenu(id uint, req dto.UpdateMenuRequest) (*models.Menu, error) {
	menu, err := s.menuRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("menu tidak ditemukan")
	}
	// (Opsional) validasi parent
	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, fmt.Errorf("parent_id tidak boleh sama dengan id menu")
		}

		// cegah siklus
		if ok, err := s.isDescendant(menu.ID, *req.ParentID); err != nil {
			return nil, err
		} else if ok {
			return nil, fmt.Errorf("parent_id mengakibatkan siklus hirarki")
		}

		if _, err := s.menuRepository.GetByID(*req.ParentID); err != nil {
			return nil, fmt.Errorf("parent menu tidak ditemukan")
		}
	}

	menu.Nama = req.Name
	menu.URL = &req.URL
	menu.Icon = req.Icon
	menu.IsActive = req.IsActive
	menu.ParentID = req.ParentID

	if err := s.menuRepository.Update(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *MenuService) GetAllMenu() ([]models.Menu, error){
	menu, err := s.menuRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return menu, nil
}


// func (s *MenuService) GetMenuByRole(roleID uint, userID string) ([]models.Menu, error) {
// 	var menus []models.Menu

// 	err := s.DB.
// 		Joins("JOIN ss_m_akses_menu ON ss_m_akses_menu.id_menu = ss_m_menu.id_menu").
// 		Where("ss_m_akses_menu.id_role = ? AND ss_m_akses_menu.can_view = ? AND ss_m_menu.is_active = ?", roleID, true, true).
// 		Order("COALESCE(parent_id, 0), id_menu").
// 		Find(&menus).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Hanya filter khusus untuk student
// 	if roleID == constants.RoleStudent {
// 		assignment, err := s.assignmentRepo.GetActiveByUserID(userID)
// 		if err != nil || assignment == nil {
// 			menus = filterMenusByURL(menus, "application/progress_report")
// 		}
// 	}

// 	return menus, nil
// }

func (s *MenuService) GetMenuByRole(roleID uint, userID string) ([]models.Menu, error) {
    var menus []models.Menu
    err := s.DB.
        Model(&models.Menu{}).
        Select("DISTINCT ss_m_menu.*"). // <-- cegah duplikat
        Joins("JOIN ss_m_akses_menu ON ss_m_akses_menu.id_menu = ss_m_menu.id_menu").
        Where("ss_m_akses_menu.id_role = ? AND ss_m_akses_menu.can_view = ? AND ss_m_menu.is_active = ?",
            roleID, true, true).
        Order("COALESCE(parent_id, 0), id_menu").
        Find(&menus).Error
    if err != nil { return nil, err }

    if roleID == constants.RoleStudent {
        assignment, err := s.assignmentRepo.GetActiveByUserID(userID)
        if err != nil || assignment == nil {
            menus = filterMenusByURL(menus, "application/progress_report")
        }
    }
    return menus, nil
}

func filterMenusByURL(menus []models.Menu, excludedURL string) []models.Menu {
	filtered := make([]models.Menu, 0)

	for _, m := range menus {
		if m.URL != nil && *m.URL == excludedURL {
			continue
		}
		filtered = append(filtered, m)
	}

	return filtered
}

func (s *MenuService) DeleteMenu(id uint) error {
    // (Opsional) tolak delete jika punya child
    var cnt int64
    if err := s.DB.Model(&models.Menu{}).
        Where("parent_id = ?", id).
        Count(&cnt).Error; err != nil {
        return err
    }
    if cnt > 0 {
        return fmt.Errorf("menu memiliki sub-menu, hapus/relokasi child terlebih dahulu")
    }
    return s.menuRepository.Delete(id)
}

func (s *MenuService) isDescendant(ancestorID, childID uint) (bool, error) {
    current := childID
    for {
        var parentID *uint
        if err := s.DB.Model(&models.Menu{}).
            Select("parent_id").
            Where("id_menu = ?", current).
            Scan(&parentID).Error; err != nil {
            return false, err
        }
        if parentID == nil {
            return false, nil
        }
        if *parentID == ancestorID {
            return true, nil
        }
        current = *parentID
    }
}