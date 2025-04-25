package dto

type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
}

type TagDTO struct {
	Name string `json:"name"`
}