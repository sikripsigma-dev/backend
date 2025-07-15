package service

import (
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
