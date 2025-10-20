package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudyProgramService interface {
	GetAll(ctx context.Context) ([]models.StudyProgram, error)
	GetByID(ctx context.Context, id string) (*models.StudyProgram, error)
	Create(ctx context.Context, name string) (*models.StudyProgram, error)
	Update(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}

type studyProgramService struct {
	repo repository.StudyProgramRepository
}

func NewStudyProgramService(repo repository.StudyProgramRepository) StudyProgramService {
	return &studyProgramService{repo}
}

func (s *studyProgramService) GetAll(ctx context.Context) ([]models.StudyProgram, error) {
	return s.repo.FindAll(ctx)
}

func (s *studyProgramService) GetByID(ctx context.Context, id string) (*models.StudyProgram, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *studyProgramService) Create(ctx context.Context, name string) (*models.StudyProgram, error) {
	item := &models.StudyProgram{
		ID:   uuid.New().String(),
		Name: name,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, fmt.Errorf("nama program studi '%s' sudah digunakan", name)
		}
		return nil, err
	}
	return item, nil
}

func (s *studyProgramService) Update(ctx context.Context, id string, name string) error {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("program studi dengan ID %s tidak ditemukan", id)
		}
		return err
	}
	item.Name = name
	item.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, item); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			return fmt.Errorf("nama program studi '%s' sudah digunakan", name)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("program studi dengan ID %s tidak ditemukan", id)
		}
		return err
	}
	return nil
}

func (s *studyProgramService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("program studi dengan ID %s tidak ditemukan", id)
		}
		return err
	}
	return nil
}
