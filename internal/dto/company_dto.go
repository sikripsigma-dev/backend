package dto

type CreateCompanyRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Description string `json:"description"`
}

type UpdateCompanyRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Description string `json:"description"`
}
