package dto

type CreateCompanyRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Description string `json:"description"`
	Industry    string `json:"industry" validate:"required"`
	Website     string `json:"website"`
	Status      string `json:"status" validate:"oneof=active inactive"` // optional kalau default, tapi disediakan di form
	Logo        string `json:"logo"`
}

type UpdateCompanyRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Description string `json:"description"`
	Industry    string `json:"industry" validate:"required"`
	Website     string `json:"website"`
	Status      string `json:"status" validate:"oneof=active inactive"`
	Logo        string `json:"logo"`
}
