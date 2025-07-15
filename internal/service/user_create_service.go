package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type UserCreateService interface {
	Create(user *models.User) error
	CreateSupervisor(sup *models.SupervisorUser) error
	CreateCompanyUser(comp *models.CompanyUser) error
	CreateUserFull(user *models.User, sup *models.SupervisorUser, comp *models.CompanyUser) error
}

type userCreateService struct {
	repo repository.UserCreateRepository
}

func NewUserCreateService(repo repository.UserCreateRepository) UserCreateService {
	return &userCreateService{repo}
}

func (s *userCreateService) Create(user *models.User) error {
	return s.repo.Create(user)
}

func (s *userCreateService) CreateSupervisor(sup *models.SupervisorUser) error {
	return s.repo.CreateSupervisor(sup)
}

func (s *userCreateService) CreateCompanyUser(comp *models.CompanyUser) error {
	return s.repo.CreateCompanyUser(comp)
}

func (s *userCreateService) CreateUserFull(
	user *models.User,
	supervisor *models.SupervisorUser,
	companyUser *models.CompanyUser,
) error {
	return s.repo.CreateUserWithRelated(user, supervisor, companyUser)
}
