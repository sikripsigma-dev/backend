package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	// "Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/util"
	"fmt"

	"mime/multipart"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type UserService struct {
	userRepository repository.UserRepository
	db             *gorm.DB
}

func NewUserService(userRepository repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{userRepository, db}
}

func (s *UserService) UpdateStudentProfile(userID string, req dto.UpdateStudentProfileRequest) error {
	// Ambil user lengkap dengan relasi
	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Pastikan student relasi tersedia
	if user.Student == nil {
		return fmt.Errorf("student profile not found")
	}

	// Mulai transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update tabel users
		user.Name = req.Name
		user.Email = req.Email
		user.Phone = req.Phone

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Update tabel student_users
		student := user.Student
		student.Nim = req.Nim
		student.Jurusan = req.Jurusan
		student.Gpa = req.Gpa
		student.UniversityID = req.UniversityID
		student.LinkedIn = req.LinkedIn
		student.Description = req.Description

		if err := tx.Save(student).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *UserService) UpdateUserCompanyProfile(userID string, req dto.UpdateUserCompanyProfileRequest) error {
	// Ambil user + relasi
	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.Company == nil {
		return fmt.Errorf("company profile not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// update tabel user
		user.Name = req.Name
		user.Email = req.Email
		user.Phone = req.Phone

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// update tabel company_user
		company := user.Company
		company.Division = req.Division
		// company.CompanyID = req.CompanyID

		if err := tx.Save(company).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *UserService) UpdateSupervisorProfile(userID string, req dto.UpdateSupervisorProfileRequest) error {
	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.Supervisor == nil {
		return fmt.Errorf("supervisor profile not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		user.Name = req.Name
		user.Email = req.Email
		user.Phone = req.Phone

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		supervisor := user.Supervisor
		supervisor.Nidn = req.NIDN

		if err := tx.Save(supervisor).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *UserService) UpdateAdminProfile(userID string, req dto.UpdateAdminProfileRequest) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone

	return s.userRepository.Update(user)
}

func (s *UserService) UpdateProfilePhoto(userID string, file *multipart.FileHeader) (string, error) {
	dir := "./public/images/user"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", userID, ext)
	filePath := filepath.Join(dir, fileName)

	// Hapus file lama
	pattern := filepath.Join(dir, fmt.Sprintf("%s.*", userID))
	matches, _ := filepath.Glob(pattern)
	for _, match := range matches {
		os.Remove(match)
	}

	// Simpan file baru
	if err := util.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Update DB pakai URL untuk FE
	imageURL := "/images/user/" + fileName
	if err := s.userRepository.UpdateImage(userID, imageURL); err != nil {
		return "", fmt.Errorf("failed to update DB: %w", err)
	}

	return imageURL, nil
}

func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, dto.UserResponse{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
			Role:  u.RoleId,
			Status: u.Status,
			Image: u.Image,
		})
	}

	return result, nil
}

func (s *UserService) UpdateUser(userID string, req dto.UpdateUserRequest) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone
	user.Status = req.Status

	return s.userRepository.Update(user)
}

func (s *UserService) GetStudentDetailByID(userID string) (*models.User, error) {
	return s.userRepository.GetStudentDetailByID(userID)
}

