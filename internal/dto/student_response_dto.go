package dto

type StudentResponse struct {
	Id                string  `json:"id"`
	Name              string  `json:"name"`
	Nim               string  `json:"nim"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	Image             string  `json:"avatar"`
	ResearchCaseTitle *string `json:"research_case_title"`
}
