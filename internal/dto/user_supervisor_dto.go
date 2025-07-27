package dto

type UpdateSupervisorProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	NIDN  string `json:"nidn"`
}
