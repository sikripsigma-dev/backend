package repository

import (
	"Skripsigma-BE/internal/models"
	"log"

	"gorm.io/gorm"
)

type UserCreateRepository interface {
	Create(user *models.User) error
	CreateSupervisor(sup *models.SupervisorUser) error
	CreateCompanyUser(comp *models.CompanyUser) error
	CreateUserWithRelated(user *models.User, sup *models.SupervisorUser, comp *models.CompanyUser) error
}

type userCreateRepository struct {
	db *gorm.DB
}

func NewUserCreateRepository(db *gorm.DB) UserCreateRepository {
	return &userCreateRepository{db}
}

func (r *userCreateRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userCreateRepository) CreateSupervisor(sup *models.SupervisorUser) error {
	return r.db.Create(sup).Error
}

func (r *userCreateRepository) CreateCompanyUser(comp *models.CompanyUser) error {
	return r.db.Create(comp).Error
}

func (r *userCreateRepository) CreateUserWithRelated(
	user *models.User,
	supervisor *models.SupervisorUser,
	companyUser *models.CompanyUser,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		log.Println("🧾 Creating user:", user.Id)

		if err := tx.Create(user).Error; err != nil {
			log.Println("❌ Failed to create user:", err)
			return err
		}
		if supervisor != nil {
			log.Println("🧾 Creating supervisor:", supervisor.UserID)
			if err := tx.Create(supervisor).Error; err != nil {
				log.Println("❌ Failed to create supervisor:", err)
				return err
			}
		}
		if companyUser != nil {
			log.Println("🧾 Creating company user:", companyUser.UserID)
			if err := tx.Create(companyUser).Error; err != nil {
				log.Println("❌ Failed to create company user:", err)
				return err
			}
		}
		return nil
	})
}
