package dto

// type RegisterRequest struct {
// 	Name            string `json:"name"`
// 	Email           string `json:"email"`
// 	Password        string `json:"password"`
// 	Phone           string `json:"phone"`
// 	SupervisorEmail string `json:"supervisor_email"`

// 	UniversityID string  `json:"university_id"`
// 	Jurusan      string  `json:"jurusan"`
// 	Gpa          float64 `json:"gpa"`
// 	Nim          string  `json:"nim"`
// }

type RegisterRequest struct {
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	Phone          string  `json:"phone"`
	UniversityID   string  `json:"university_id"`
	StudyProgramID string  `json:"study_program_id"`
	Jurusan        string  `json:"jurusan"`
	Gpa            float64 `json:"gpa"`
	Nim            string  `json:"nim"`

	// (opsional/backward compat—tak dipakai lagi untuk mapping)
	SupervisorEmail string `json:"supervisor_email,omitempty"`
}


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
