package dto

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type RoleDTO struct {
	Name string `json:"name"`
}