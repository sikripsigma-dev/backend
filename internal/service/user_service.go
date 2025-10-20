package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/util"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type UserService struct {
	userRepository repository.UserRepository
	db             *gorm.DB
}

func NewUserService(userRepository repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{userRepository, db}
}

// ==============================
// Student Profile
// ==============================
func (s *UserService) UpdateStudentProfile(userID string, req dto.UpdateStudentProfileRequest) error {
	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if user.Student == nil {
		return fmt.Errorf("student profile not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// users
		user.Name = req.Name
		user.Email = req.Email
		user.Phone = req.Phone
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// student_users
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

// ==============================
// Company Profile
// ==============================
func (s *UserService) UpdateUserCompanyProfile(userID string, req dto.UpdateUserCompanyProfileRequest) error {
	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if user.Company == nil {
		return fmt.Errorf("company profile not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// users
		user.Name = req.Name
		user.Email = req.Email
		user.Phone = req.Phone
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// company_users
		company := user.Company
		company.Division = req.Division
		if err := tx.Save(company).Error; err != nil {
			return err
		}
		return nil
	})
}

// ==============================
// Supervisor Profile
// ==============================
func (s *UserService) UpdateSupervisorProfile(userID string, req dto.UpdateSupervisorProfileRequest) error {
	user, err := s.userRepository.GetWithRelationsByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if user.Supervisor == nil {
		return fmt.Errorf("supervisor profile not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// users
		user.Name = req.Name
		user.Email = req.Email
		user.Phone = req.Phone
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// supervisor_users
		supervisor := user.Supervisor
		supervisor.Nidn = req.NIDN
		if err := tx.Save(supervisor).Error; err != nil {
			return err
		}
		return nil
	})
}

// ==============================
// Admin Profile (umum)
// ==============================
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

// ==============================
// Profile Photo
// ==============================
func (s *UserService) UpdateProfilePhoto(userID string, file *multipart.FileHeader) (string, error) {
	dir := "./public/images/user"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", userID, ext)
	filePath := filepath.Join(dir, fileName)

	// hapus file lama pola {userID}.*
	pattern := filepath.Join(dir, fmt.Sprintf("%s.*", userID))
	if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
		for _, m := range matches {
			_ = os.Remove(m)
		}
	}

	if err := util.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	imageURL := "/images/user/" + fileName
	if err := s.userRepository.UpdateImage(userID, imageURL); err != nil {
		return "", fmt.Errorf("failed to update DB: %w", err)
	}
	return imageURL, nil
}

// ==============================
// Get All Users
// ==============================
func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, err
	}
	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, dto.UserResponse{
			Id:     u.Id,
			Name:   u.Name,
			Email:  u.Email,
			Phone:  u.Phone,
			Role:   u.RoleId,
			Status: u.Status,
			Image:  u.Image,
		})
	}
	return result, nil
}

// ==============================
// Update User (ADMIN) + Upsert Kaprodi.Prodi
// ==============================
//
// - Update data dasar user.
// - Jika request menyertakan study_program_id DAN user adalah Kaprodi (role_id = 5),
//   maka lakukan upsert ke tabel headstudy_users (update/insert) dan sinkron UniversityID
//   dari tabel study_programs.
//
// Update User (ADMIN) + Upsert Kaprodi.Prodi (tanpa UniversityID)
func (s *UserService) UpdateUser(userID string, req dto.UpdateUserRequest) error {
    user, err := s.userRepository.GetByID(userID)
    if err != nil {
        return fmt.Errorf("user not found")
    }

    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1) update users
        user.Name = req.Name
        user.Email = req.Email
        user.Phone = req.Phone
        if req.Status != "" {
            user.Status = req.Status
        }
        if err := tx.Save(user).Error; err != nil {
            return err
        }

        // 2) kalau user ini Kaprodi (role_id=5), kelola mapping prodi+universitas
        if user.RoleId == 5 {
            // wajib keduanya disediakan saat assign
            if req.StudyProgramID == nil || req.UniversityID == nil {
                // Boleh dibuat optional: hanya update kalau ada isian.
                // Kalau ingin strict, return error:
                return fmt.Errorf("study_program_id and university_id are required for Kaprodi")
            }

            spID := strings.TrimSpace(*req.StudyProgramID)
            univID := strings.TrimSpace(*req.UniversityID)

            // validasi keberadaan master
            if spID != "" {
                var sp models.StudyProgram
                if err := tx.First(&sp, "id = ?", spID).Error; err != nil {
                    return fmt.Errorf("study program not found")
                }
            }
            if univID != "" {
                var univ models.University
                if err := tx.First(&univ, "id = ?", univID).Error; err != nil {
                    return fmt.Errorf("university not found")
                }
            }

            // upsert ss_t_headstudy_users by user_id
			var hs models.HeadstudyUser
			err := tx.Where("user_id = ?", user.Id).First(&hs).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// ⇢ CREATE (wajib isi NIDN kalau kolomnya NOT NULL)
				if req.Nidn == nil || strings.TrimSpace(*req.Nidn) == "" {
					return fmt.Errorf("nidn is required for creating headstudy_user")
				}
				hs = models.HeadstudyUser{
					UserID:         user.Id,           // PK
					UniversityID:   univID,
					StudyProgramID: spID,
					Nidn:           strings.TrimSpace(*req.Nidn),
				}
				if err := tx.Create(&hs).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				// ⇢ UPDATE (optional: hanya update jika field dikirim)
				if req.UniversityID != nil {
					hs.UniversityID = univID
				}
				if req.StudyProgramID != nil {
					hs.StudyProgramID = spID
				}
				if req.Nidn != nil && strings.TrimSpace(*req.Nidn) != "" {
					hs.Nidn = strings.TrimSpace(*req.Nidn)
				}
				if err := tx.Save(&hs).Error; err != nil {
					return err
				}
			}

        }

        return nil
    })
}

func (s *UserService) GetUserWithRelations(userID string) (*models.User, error) {
	return s.userRepository.GetWithRelationsByID(userID)
}

// ==============================
// Others
// ==============================
func (s *UserService) GetStudentDetailByID(userID string) (*models.User, error) {
	return s.userRepository.GetStudentDetailByID(userID)
}

func (s *UserService) AssignSupervisorToStudent(studentID, supervisorID string) error {
	return s.userRepository.AssignSupervisorToStudent(studentID, supervisorID)
}
