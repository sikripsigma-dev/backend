package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type AuthTokenRepository interface {
	Create(token *models.AuthToken) error
	MarkUsed(id string) error
	GetByToken(token string) (*models.AuthToken, error)
	DeleteByID(id string) error
}

type authTokenRepo struct {
	db *gorm.DB
}

func NewAuthTokenRepository(db *gorm.DB) AuthTokenRepository {
	return &authTokenRepo{db}
}

func (r *authTokenRepo) Create(token *models.AuthToken) error {
	return r.db.Create(token).Error
}

func (r *authTokenRepo) GetByToken(token string) (*models.AuthToken, error) {
	var authToken models.AuthToken
	if err := r.db.Where("token = ?", token).First(&authToken).Error; err != nil {
		return nil, err
	}
	return &authToken, nil
}

func (r *authTokenRepo) DeleteByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.AuthToken{}).Error
}


func (r *authTokenRepo) MarkUsed(id string) error {
	return r.db.Model(&models.AuthToken{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}
