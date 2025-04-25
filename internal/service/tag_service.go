package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type TagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) *TagService {
	return &TagService{tagRepo}
}

func (s *TagService) CreateTag(req dto.CreateTagRequest) (*models.Tag, error) {
	
	tag := &models.Tag{
		Name:        req.Name,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) GetTagByID(id string) (*models.Tag, error) {
	
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (s* TagService) GetAllTags() ([]models.Tag, error) {
	
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return tags, nil
}