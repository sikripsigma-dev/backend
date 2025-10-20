package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/util"
	"fmt"
)

type MonitoringProgressService interface {
	GiveFeedback(payload dto.CreateSupervisorFeedbackDTO) error
	GiveFeedbackCompany(payload dto.CreateCompanyFeedbackDTO) error
	GetSupervisorFeedback(reportID uint) ([]models.SupervisorMonitoringProgress, error)
	GetCompanyFeedback(reportID uint) ([]models.CompanyMonitoringProgress, error)
}

type monitoringProgressService struct {
	repo                repository.MonitoringProgressRepository
	weeklyReportRepo    repository.WeeklyReportRepository
	companyWeeklyRepo   repository.CompanyWeeklyReportRepository
	notificationService NotificationService
	userRepository      repository.UserRepository
}

func NewMonitoringProgressService(
	r repository.MonitoringProgressRepository,
	weeklyReportRepo repository.WeeklyReportRepository,
	companyWeeklyRepo repository.CompanyWeeklyReportRepository,
	notificationService NotificationService,
	userRepository repository.UserRepository,
) MonitoringProgressService {
	return &monitoringProgressService{
		repo:                r,
		weeklyReportRepo:    weeklyReportRepo,
		companyWeeklyRepo:   companyWeeklyRepo,
		notificationService: notificationService,
		userRepository:      userRepository,
	}
}

// ============ Feedback oleh Supervisor ============
func (s *monitoringProgressService) GiveFeedback(payload dto.CreateSupervisorFeedbackDTO) error {
	// 1. Simpan feedback dosen
	data := models.SupervisorMonitoringProgress{
		WeeklyReportID: payload.WeeklyReportID,
		Feedback:       payload.Feedback,
		Status:         payload.Status,
	}
	if err := s.repo.Create(&data); err != nil {
		return err
	}

	// 2. Update status laporan mingguan
	if err := s.repo.UpdateReportStatus(payload.WeeklyReportID, payload.Status); err != nil {
		return err
	}

	// 3. Ambil informasi laporan mingguan
	report, err := s.weeklyReportRepo.GetByID(payload.WeeklyReportID)
	if err != nil {
		return fmt.Errorf("Gagal mengambil laporan mingguan: %v", err)
	}

	// 4. Kirim notifikasi ke mahasiswa
	notif := &models.Notification{
		UserID:  &report.StudentID,
		Type:    "weekly_report_reviewed",
		Message: fmt.Sprintf("Laporan Minggu ke-%d Anda telah ditinjau oleh dosen pembimbing", report.Week),
		Metadata: util.JSONB{
			"weekly_report_id": report.ID,
			"status":           payload.Status,
			"week":             report.Week,
			"student_id":       report.StudentID,
		},
	}
	_ = s.notificationService.CreateNotification(notif) // abaikan error

	return nil
}

// ============ Feedback oleh Perusahaan ============
func (s *monitoringProgressService) GiveFeedbackCompany(payload dto.CreateCompanyFeedbackDTO) error {
	// 1. Simpan feedback perusahaan
	data := models.CompanyMonitoringProgress{
		WeeklyReportID: payload.WeeklyReportID,
		Feedback:       payload.Feedback,
	}
	if err := s.repo.CompanyCreate(&data); err != nil {
		return err
	}

	// 2. Ambil laporan mingguan dari perusahaan
	report, err := s.companyWeeklyRepo.GetByID(payload.WeeklyReportID)
	if err != nil {
		return fmt.Errorf("Gagal mengambil laporan mingguan perusahaan: %v", err)
	}

	// 3. Ambil data user (mahasiswa)
	user, err := s.userRepository.GetByID(report.StudentID)
	if err != nil {
		return fmt.Errorf("Gagal mengambil data mahasiswa: %v", err)
	}

	// 4. Kirim notifikasi ke mahasiswa
	notif := &models.Notification{
		UserID:  &report.StudentID,
		Type:    "company_feedback",
		Message: fmt.Sprintf("Perusahaan telah memberikan feedback untuk laporan Minggu ke-%d", report.Week),
		Metadata: util.JSONB{
			"report_id":         report.ID,
			"week":              report.Week,
			"research_case_id":  report.ResearchCaseID,
			"student_name":      user.Name,
			"feedback_from":     "company",
		},
	}
	_ = s.notificationService.CreateNotification(notif) // abaikan error

	return nil
}

func (s *monitoringProgressService) GetSupervisorFeedback(reportID uint) ([]models.SupervisorMonitoringProgress, error) {
	return s.repo.GetSupervisorFeedbackByReportID(reportID)
}

func (s *monitoringProgressService) GetCompanyFeedback(reportID uint) ([]models.CompanyMonitoringProgress, error) {
	return s.repo.GetCompanyFeedbackByReportID(reportID)
}

