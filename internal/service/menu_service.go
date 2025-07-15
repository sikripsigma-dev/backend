package service

import (
	"Skripsigma-BE/internal/constants"
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

// type MenuService struct {
// 	menuRepository repository.MenuRepository
// 	DB             *gorm.DB
// }

// func NewMenuService(menuRepository repository.MenuRepository, db *gorm.DB) *MenuService {
// 	return &MenuService{
// 		menuRepository: menuRepository,
// 		DB:             db,
// 	}
// }

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

func (s *MenuService) CreateMenu(req dto.CreateMenuRequest) (*models.Menu, error) {
	menu := &models.Menu{
		Nama:    	req.Name,
		URL:         &req.URL,
		ParentID: req.ParentID,
	}

	if err := s.menuRepository.Create(menu); err != nil {
		return nil, err
	}
	return menu, nil

}

func (s *MenuService) UpdateMenu(ID string, req dto.UpdateMenuRequest) (*models.Menu, error) {
	menu, err := s.menuRepository.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("Menu Not Found")
	}

	menu.Nama = req.Name
	menu.URL = &req.URL
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

// func (s *MenuService) GetMenuByRole(RoleID uint) ([]models.Menu, error){
// 	var menus []models.Menu

// 	err := s.DB.
// 		Joins("JOIN ss_m_akses_menu ON ss_m_akses_menu.id_menu = ss_m_menu.id_menu").
// 			Where("ss_m_akses_menu.id_role = ? AND ss_m_akses_menu.can_view = ?", RoleID, true).
// 			Find(&menus).Error

// 	return menus, err
// }

func (s *MenuService) GetMenuByRole(roleID uint, userID string) ([]models.Menu, error) {
	var menus []models.Menu

	err := s.DB.
		Joins("JOIN ss_m_akses_menu ON ss_m_akses_menu.id_menu = ss_m_menu.id_menu").
		Where("ss_m_akses_menu.id_role = ? AND ss_m_akses_menu.can_view = ?", roleID, true).
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	// Hanya filter khusus untuk student
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