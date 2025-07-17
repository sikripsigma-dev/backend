package service

import (
	"Skripsigma-BE/internal/dto"
	// "Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/util"
	"fmt"

	"mime/multipart"
	"os"
	"path/filepath"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository}
}

func (s *UserService) UpdateStudentProfile(userID string, req dto.UpdateStudentProfileRequest) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone

	return s.userRepository.Update(user)
}

func (s *UserService) UpdateUserCompanyProfile(userID string, req dto.UpdateUserCompanyProfileRequest) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone

	return s.userRepository.Update(user)
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

