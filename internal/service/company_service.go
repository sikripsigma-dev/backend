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

// func (s *CompanyService) CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error) {
	
// 	company := &models.Company{
// 		Name:        req.Name,
// 		Description: req.Description,
// 		Email:       req.Email,
// 		Phone:	   	 req.Phone,
// 		Address:     req.Address,
// 	}

// 	if err := s.companyRepo.Create(company); err != nil {
// 		return nil, err
// 	}
// 	return company, nil
// }

func (s *CompanyService) CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error) {
	company := &models.Company{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		Description: req.Description,
		Industry:    req.Industry,
		Website:     req.Website,
		Status:      req.Status,
		Logo:        req.Logo,
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

// func (s *CompanyService) UpdateCompany(id string, req dto.UpdateCompanyRequest) (*models.Company, error) {
// 	company, err := s.companyRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	company.Name = req.Name
// 	company.Email = req.Email
// 	company.Phone = req.Phone
// 	company.Address = req.Address
// 	company.Description = req.Description

// 	if err := s.companyRepo.Update(company); err != nil {
// 		return nil, err
// 	}
// 	return company, nil
// }

func (s *CompanyService) UpdateCompany(id string, req dto.UpdateCompanyRequest) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	company.Name = req.Name
	company.Email = req.Email
	company.Phone = req.Phone
	company.Address = req.Address
	company.Description = req.Description
	company.Industry = req.Industry
	company.Website = req.Website
	company.Status = req.Status
	company.Logo = req.Logo

	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}
	return company, nil
}