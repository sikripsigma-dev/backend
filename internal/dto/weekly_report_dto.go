// dto/weekly_report_dto.go
package dto

type CreateWeeklyReportDTO struct {
	ResearchCaseID string `validate:"required,uuid4"`
	Week      int      `form:"week" validate:"required"`
	Progress  string   `form:"progress" validate:"required"`
	Plans     string   `form:"plans"`
	Mood      int      `form:"mood" validate:"required"`
	Notes     string   `form:"notes"`
	StartDate string   `form:"start_date"`
	EndDate   string   `form:"end_date"`
	Files     []string `form:"-"` // tidak di-parse langsung, ditangani manual
}


type CreateCompanyWeeklyReportDTO struct {
	ResearchCaseID string   `form:"research_case_id" validate:"required,uuid4"`
	Week           int      `form:"week" validate:"required,min=1,max=16"`
	Activities     string   `form:"activities" validate:"required"`
	Issues         string   `form:"issues"`
	Hopes          string   `form:"hopes"`
	Notes          string   `form:"notes"`
	StartDate      string   `form:"start_date"`
	EndDate        string   `form:"end_date"`
	Files          []string `form:"files"`
}


type WeeklyReportResponse struct {
	ID        uint     `json:"id"`
	Week      int      `json:"week"`
	Progress  string   `json:"progress"`
	Plans     string   `json:"plans"`
	Mood      int      `json:"mood"`
	Notes     string   `json:"notes"`
	Status    string   `json:"status"`
	StartDate string   `json:"start_date"` 
	EndDate   string   `json:"end_date"` 
	Files     []string `json:"files"`
}

type CompanyWeeklyReportResponse struct {
	ID        uint     `json:"id"`
	Week      int      `json:"week"`
	Activities string   `json:"activities"`
	Issues     string   `json:"issues"`
	Hopes      string   `json:"hopes"`
	Notes      string   `json:"notes"`
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
	Files      []string `json:"files"`

	Student struct {
		Name     string `json:"name"`
		NIM      string `json:"nim"`
		Avatar   string `json:"avatar"`
		Initials string `json:"initials"`
	} `json:"student"`

	ResearchCaseTitle string `json:"research_case_title"`
}


type SimpleStudentDTO struct {
	Name     string `json:"name"`
	NIM      string `json:"nim"`
	Avatar   string `json:"avatar"`
	Initials string `json:"initials"`
}
