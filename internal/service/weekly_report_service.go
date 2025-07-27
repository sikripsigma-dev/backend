package service

import (
	"strings"
	"time"

	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"

	"encoding/json"
	"fmt"
)

type WeeklyReportService interface {
	SubmitReport(studentID string, dto dto.CreateWeeklyReportDTO) error
	GetReports(studentID string) ([]dto.WeeklyReportResponse, error)
	GetReportsBySupervisor(universityID string) ([]models.WeeklyReport, error)
	SubmitCompanyReport(studentID string, dto dto.CreateCompanyWeeklyReportDTO) error
	GetReportsByCompany(companyID string) ([]dto.CompanyWeeklyReportResponse, error)
}

type weeklyReportService struct {
	repo repository.WeeklyReportRepository
}

func NewWeeklyReportService(r repository.WeeklyReportRepository) WeeklyReportService {
	return &weeklyReportService{r}
}


func (s *weeklyReportService) SubmitReport(studentID string, input dto.CreateWeeklyReportDTO) error {
	// Format tanggal yang dikirim dari frontend
	const layout = "2006-01-02"

	// Parse string tanggal ke time.Time
	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		return fmt.Errorf("format start_date tidak valid: %w", err)
	}

	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		return fmt.Errorf("format end_date tidak valid: %w", err)
	}

	// Marshal file names ke JSON string
	filesJSON, err := json.Marshal(input.Files)
	if err != nil {
		return fmt.Errorf("gagal memproses file: %w", err)
	}

	// Bangun model WeeklyReport
	report := models.WeeklyReport{
		StudentID: studentID,
		ResearchCaseID: input.ResearchCaseID,
		Week:      input.Week,
		Progress:  input.Progress,
		Plans:     input.Plans,
		Mood:      input.Mood,
		Notes:     input.Notes,
		Status:    "Menunggu Review",
		StartDate: startDate,
		EndDate:   endDate,
		Files:     string(filesJSON),
	}

	return s.repo.Create(&report)
}

func (s *weeklyReportService) SubmitCompanyReport(studentID string, input dto.CreateCompanyWeeklyReportDTO) error {
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return fmt.Errorf("format start_date tidak valid: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return fmt.Errorf("format end_date tidak valid: %w", err)
	}

	filesJSON, err := json.Marshal(input.Files)
	if err != nil {
		return fmt.Errorf("gagal memproses file: %w", err)
	}

	report := models.CompanyWeeklyReport{
		StudentID:      studentID,
		ResearchCaseID: input.ResearchCaseID,
		Week:           input.Week,
		Activities:     input.Activities,
		Issues:         input.Issues,
		Hopes:          input.Hopes,
		Notes:          input.Notes,
		StartDate:      startDate,
		EndDate:        endDate,
		Files:          string(filesJSON),
	}

	return s.repo.CreateCompanyReport(&report)
}


func (s *weeklyReportService) GetReports(studentID string) ([]dto.WeeklyReportResponse, error) {
	reports, err := s.repo.GetByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var response []dto.WeeklyReportResponse
	for _, r := range reports {
		var files []string
		if err := json.Unmarshal([]byte(r.Files), &files); err != nil {
			files = []string{}
		}
		response = append(response, dto.WeeklyReportResponse{
			ID:        r.ID,
			Week:      r.Week,
			Progress:  r.Progress,
			Plans:     r.Plans,
			Mood:      r.Mood,
			Notes:     r.Notes,
			Status:    r.Status,
			StartDate: r.StartDate.Format("02-01-2006"),
			EndDate:   r.EndDate.Format("02-01-2006"),
			Files:     files,
		})
	}

	return response, nil
}

func (s *weeklyReportService) GetReportsBySupervisor(universityID string) ([]models.WeeklyReport, error) {
	return s.repo.GetByUniversity(universityID)
}

func (s *weeklyReportService) GetReportsByCompany(companyID string) ([]dto.CompanyWeeklyReportResponse, error) {
	reports, err := s.repo.GetByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	var response []dto.CompanyWeeklyReportResponse
	for _, r := range reports {
		var files []string
		if err := json.Unmarshal([]byte(r.Files), &files); err != nil {
			files = []string{}
		}

		// Generate inisial mahasiswa
		initial := ""
		nameParts := strings.Split(r.Student.Name, " ")
		for _, part := range nameParts {
			if len(part) > 0 {
				initial += strings.ToUpper(part[:1])
			}
		}

		response = append(response, dto.CompanyWeeklyReportResponse{
			ID:         r.ID,
			Week:       r.Week,
			Activities: r.Activities,
			Issues:     r.Issues,
			Hopes:      r.Hopes,
			Notes:      r.Notes,
			StartDate:  r.StartDate.Format("02-01-2006"),
			EndDate:    r.EndDate.Format("02-01-2006"),
			Files:      files,
			Student: dto.SimpleStudentDTO{
				Name:     r.Student.Name,
				NIM:      r.Student.Nim,
				Avatar:   r.Student.Image,
				Initials: initial,
			},
			ResearchCaseTitle: r.ResearchCase.Title,
		})
	}

	return response, nil
}

