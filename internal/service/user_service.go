package service

import (
	"Skripsigma-BE/internal/dto"
	// "Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
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


// func (s *UserService) UpdateStudentProfile(userID string, req dto.UpdateStudentProfileRequest) (*models.User, error) {
// 	user, err := s.userRepository.GetByID(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("User not found")
// 	}

// 	// Update user fields with the request data
// 	user.Name = req.Name
// 	user.Email = req.Email
// 	user.Phone = req.Phone

// 	if err := s.userRepository.Update(user); err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }