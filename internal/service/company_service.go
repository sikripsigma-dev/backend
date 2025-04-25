package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type CompanyService struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo}
}

func (s *CompanyService) CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error) {
	
	company := &models.Company{
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		Phone:	   	 req.Phone,
		Address:     req.Address,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}
	return company, nil
}

func (s *CompanyService) GetCompanyByID(id string) (*models.Company, error) {
	
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (s* CompanyService) GetAllCompanies() ([]models.Company, error) {

	comapnies, err := s.companyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return comapnies, nil
}