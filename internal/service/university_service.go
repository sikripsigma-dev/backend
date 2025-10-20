package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

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

// func (s *universityService) Create(ctx context.Context, name string) (*models.University, error) {
// 	u := &models.University{
// 		ID:   uuid.New().String(),
// 		Name: name,
// 	}
// 	err := s.repo.Create(ctx, u)
// 	return u, err
// }

func (s *universityService) Create(ctx context.Context, name string) (*models.University, error) {
	u := &models.University{ID: uuid.New().String(), Name: name}
	err := s.repo.Create(ctx, u)
	if err != nil {
		// Tangani duplikat secara portable
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
		   strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
		   strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, fmt.Errorf("nama universitas '%s' sudah digunakan", name)
		}
		return nil, err
	}
	return u, nil
}

// func (s *universityService) Update(ctx context.Context, id string, name string) error {
// 	u, err := s.repo.FindByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	u.Name = name
// 	u.UpdatedAt = time.Now()
// 	return s.repo.Update(ctx, u)
// }

func (s *universityService) Update(ctx context.Context, id string, name string) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("universitas dengan ID %s tidak ditemukan", id)
		}
		return err
	}

	u.Name = name
	u.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, u); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
		   strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
		   strings.Contains(strings.ToLower(err.Error()), "unique") {
			return fmt.Errorf("nama universitas '%s' sudah digunakan", name)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("universitas dengan ID %s tidak ditemukan", id)
		}
		return err
	}
	return nil
}

func (s *universityService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("universitas dengan ID %s tidak ditemukan", id)
	}
	return err
}