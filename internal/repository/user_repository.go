package repository

import (
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	// GetWithCompanyByID(id string) (*models.User, error)
	Update(user *models.User) error
	UpdateImage(userID string, imageURL string) error
	GetWithRelationsByID(id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (r *userRepository) GetWithCompanyByID(id string) (*models.User, error) {
// 	var user models.User
// 	if err := r.db.Preload("Company").Where("id = ?", id).First(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

func (r *userRepository) GetWithRelationsByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.
		Preload("Company").
		Preload("Student").
		Preload("Student.University").
		Preload("Supervisor").
		Preload("Supervisor.University").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateImage(userID string, imageURL string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("image", imageURL).Error
}
