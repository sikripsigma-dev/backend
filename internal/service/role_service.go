package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo}
}

func (s *RoleService) CreateRole(req dto.CreateRoleRequest) (*models.Role, error) {
	
	role := &models.Role{
		Name: req.Name,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetRoleByID(id string) (*models.Role, error) {
	
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return roles, nil
}