package dto

// type UpdateStudentProfileRequest struct {
// 	Name  string `json:"name" validate:"required"`
// 	Email string `json:"email" validate:"required,email"`
// 	Phone string `json:"phone" validate:"required"`
// }

type UpdateStudentProfileRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	Description  string  `json:"description"`
	Nim          string  `json:"nim"`
	Jurusan      string  `json:"jurusan"`
	Gpa          float64 `json:"gpa"`
	UniversityID string  `json:"university_id"`
	LinkedIn     string  `json:"linkedin"`
}