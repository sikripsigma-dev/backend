package repository

import (
	"Skripsigma-BE/internal/models"

	"time"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(application *models.Application) error
	GetByID(id string) (*models.Application, error)
	Update(application *models.Application) error
	// GetAll() ([]models.Application, error)
	GetByStudentID(studentID string) (*models.Application, error)
	GetByResearchCaseID(researchCaseID string) ([]models.Application, error)
	GetByUserIDAndResearchCaseID(userID, researchCaseID string) (*models.Application, error)
	GetAllByStudentID(studentID string) ([]models.Application, error)
	AssociateTags(researchCaseID string, tagIDs []string) error
	CancelOtherApplications(currentAppID uint, userID string) error
	GetPendingHeadstudyConfirmations(universityID, studyProgramID string) ([]models.Application, error)
	UpdateHeadStudyDecision(id uint, decision string) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db}
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(application).Error
}

// cek apabila user sudah apply
func (r *applicationRepository) HasApplied(userID, researchCaseID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Application{}).
		Where("user_id = ? AND research_case_id = ?", userID, researchCaseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *applicationRepository) GetByUserIDAndResearchCaseID(userID, researchCaseID string) (*models.Application, error) {
	var application models.Application
	err := r.db.Where("user_id = ? AND research_case_id = ?", userID, researchCaseID).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}



func (r *applicationRepository) GetByID(id string) (*models.Application, error) {
	
	var application models.Application
	if err := r.db.Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) Update(application *models.Application) error {
	return r.db.Save(application).Error
}

func (r *applicationRepository) GetByStudentID(studentID string) (*models.Application, error) {
	var application models.Application
	if err := r.db.Where("user_id = ?", studentID).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetAllByStudentID(studentID string) ([]models.Application, error) {
	var applications []models.Application
	err := r.db.
		Preload("ResearchCase").
		Preload("ResearchCase.Company").
		// .Preload("ResearchCase.Tags"). // kalau ingin menyertakan tags
		Where("user_id = ?", studentID).
		Find(&applications).Error

	if err != nil {
		return nil, err
	}
	return applications, nil
}


// get by research case id
func (r *applicationRepository) GetByResearchCaseID(researchCaseID string) ([]models.Application, error) {
	var applications []models.Application
	err := r.db.
		Preload("User").
		Preload("ResearchCase").
		Where("research_case_id = ?", researchCaseID).
		Find(&applications).Error

	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetPendingHeadstudyConfirmations(universityID, studyProgramID string) ([]models.Application, error) {
	var applications []models.Application
	err := r.db.
		Joins("JOIN ss_student_user s ON s.user_id = ss_t_applications.user_id").
		// Where("s.university_id = ? AND s.jurusan = ? AND ss_t_applications.status = ?", universityID, studyProgramID, "accepted").
		Where("s.university_id = ? AND s.study_program_id = ?", universityID, studyProgramID).
		Preload("User").
		Preload("ResearchCase").
		Preload("ResearchCase.Company").
		Find(&applications).Error

	if err != nil {
		return nil, err
	}
	return applications, nil
}


func (r *applicationRepository) UpdateHeadStudyDecision(id uint, decision string) error {
	return r.db.Model(&models.Application{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"head_study_decision": decision,
			// "processed_at":        time.Now(),
		}).Error
}


func (r *applicationRepository) AssociateTags(researchCaseID string, tagIDs []string) error {
	for _, tagID := range tagIDs {
		researchCaseTag := models.ResearchCaseTag{
			ResearchCaseID: researchCaseID,
			TagID:          tagID,
		}
		if err := r.db.Create(&researchCaseTag).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *applicationRepository) CancelOtherApplications(currentAppID uint, userID string) error {
	return r.db.
		Model(&models.Application{}).
		Where("id != ? AND user_id = ? AND status IN ?", currentAppID, userID, []string{"pending", "accepted"}).
		Updates(map[string]interface{}{
			"status":       "cancelled",
			"processed_at": time.Now(),
		}).Error
}


