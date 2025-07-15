package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type UniversityService interface {
	GetAll(ctx context.Context) ([]models.University, error)
	GetByID(ctx context.Context, id string) (*models.University, error)
	Create(ctx context.Context, name string) (*models.University, error)
	Update(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}

type universityService struct {
	repo repository.UniversityRepository
}

func NewUniversityService(repo repository.UniversityRepository) UniversityService {
	return &universityService{repo}
}

func (s *universityService) GetAll(ctx context.Context) ([]models.University, error) {
	return s.repo.FindAll(ctx)
}

func (s *universityService) GetByID(ctx context.Context, id string) (*models.University, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *universityService) Create(ctx context.Context, name string) (*models.University, error) {
	u := &models.University{
		ID:   uuid.New().String(),
		Name: name,
	}
	err := s.repo.Create(ctx, u)
	return u, err
}

func (s *universityService) Update(ctx context.Context, id string, name string) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

func (s *universityService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
