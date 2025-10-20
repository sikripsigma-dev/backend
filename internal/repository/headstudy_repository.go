package repository

import (
	"Skripsigma-BE/internal/dto"

	"gorm.io/gorm"
)

type HeadstudyRepository interface {
	GetStudentsByUniversity(univID string) ([]dto.StudentResponse, error)
}

type headstudyRepository struct {
	db *gorm.DB
}

func NewHeadstudyRepository(db *gorm.DB) HeadstudyRepository {
	return &headstudyRepository{db}
}

func (r *headstudyRepository) GetStudentsByUniversity(univID string) ([]dto.StudentResponse, error) {
	var students []dto.StudentResponse

	err := r.db.Table("ss_users").
		Select(`
			ss_users.id,
			ss_users.name,
			ss_student_user.nim,
			ss_users.email,
			ss_users.image,
			ss_users.phone,
			ss_t_research_cases.title AS research_case_title
		`).
		Joins("JOIN ss_student_user ON ss_student_user.user_id = ss_users.id").
		Joins("LEFT JOIN ss_t_assignments ON ss_t_assignments.user_id = ss_users.id AND ss_t_assignments.status = 'active'").
		Joins("LEFT JOIN ss_t_research_cases ON ss_t_research_cases.id = ss_t_assignments.research_case_id").
		Joins("JOIN ss_headstudy_user ON ss_headstudy_user.study_program_id = ss_student_user.study_program_id").
		Where("ss_headstudy_user.university_id = ?", univID).
		Scan(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}
