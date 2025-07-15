package dto

type Menu struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active"`
}

type CreateMenuRequest struct {
	Name     string `json:"name" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active" validate:"required"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateMenuRequest struct {
	Name     string `json:"name" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active" validate:"required"`
	ParentID *uint  `json:"parent_id"`
}