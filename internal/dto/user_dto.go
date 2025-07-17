package dto

type UserResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Role   uint   `json:"role_id"`
	Status string `json:"status"`
	Image  string `json:"image"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Phone  string `json:"phone" validate:"required"`
	Status string `json:"status"`
}
