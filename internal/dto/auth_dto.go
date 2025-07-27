package dto

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Phone           string `json:"phone"`
	SupervisorEmail string `json:"supervisor_email"`

	UniversityID string  `json:"university_id"`
	Jurusan      string  `json:"jurusan"`
	Gpa          float64 `json:"gpa"`
	Nim          string  `json:"nim"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
