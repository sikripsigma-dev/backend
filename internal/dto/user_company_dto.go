package dto

// type UpdateUserCompanyProfileRequest struct {
// 	Name  string `json:"name" validate:"required"`
// 	Email string `json:"email" validate:"required,email"`
// 	Phone string `json:"phone" validate:"required"`
// }

type UpdateUserCompanyProfileRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Division   string `json:"division"`
	// CompanyID  string `json:"company_id"`
}