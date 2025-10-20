package repository

import (
	"Skripsigma-BE/internal/constants"
	"Skripsigma-BE/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	CreateStudent(student *models.StudentUser) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	// GetWithCompanyByID(id string) (*models.User, error)
	Update(user *models.User) error
	UpdateImage(userID string, imageURL string) error
	GetWithRelationsByID(id string) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetStudentDetailByID(id string) (*models.User, error)
	AssignSupervisorToStudent(studentID string, supervisorID string) error
	FindHeadstudyByUnivAndProdi(univID, prodiID string) (*models.User, error)
	GetStudyProgramByID(id string) (*models.StudyProgram, error)
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


func (r *userRepository) GetWithRelationsByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.
		Preload("Company").
		Preload("Student").
		Preload("Student.University").
		Preload("Supervisor").
		Preload("Supervisor.University").
		Preload("Headstudy").
		Preload("Headstudy.University").
		Preload("Headstudy.StudyProgram").
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

func (r *userRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetStudentDetailByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.
		Preload("Student").
		Preload("Student.University").
		Preload("StudentDocuments"). // ← tambahkan relasi dokumen
		Where("id = ? AND role_id = ?", id, constants.RoleStudent).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateStudent(student *models.StudentUser) error {
	return r.db.Create(student).Error
}

func (r *userRepository) AssignSupervisorToStudent(studentID string, supervisorID string) error {
    return r.db.Model(&models.StudentUser{}).
        Where("user_id = ?", studentID).
        Update("supervisor_id", supervisorID).Error
}

func (r *userRepository) FindHeadstudyByUnivAndProdi(univID, prodiID string) (*models.User, error) {
	var u models.User
	err := r.db.
		Joins("JOIN ss_headstudy_user hs ON hs.user_id = ss_users.id").
		Where("hs.university_id = ? AND hs.study_program_id = ?", univID, prodiID).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetStudyProgramByID(id string) (*models.StudyProgram, error) {
    var sp models.StudyProgram
    if err := r.db.Where("id = ?", id).First(&sp).Error; err != nil {
        return nil, err
    }
    return &sp, nil
}


