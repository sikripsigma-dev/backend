package dto

type MenuAccessNode struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	URL      *string `json:"url,omitempty"`
	Icon     string  `json:"icon"`
	IsActive bool    `json:"is_active"`
	ParentID *uint   `json:"parent_id,omitempty"`

	CanView    bool `json:"can_view"`
	CanAdd     bool `json:"can_add"`
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
	CanExport  bool `json:"can_export"`

	Children []MenuAccessNode `json:"children,omitempty"`
}

type UpdateMenuAccessItem struct {
	MenuID     uint `json:"menu_id" validate:"required"`
	CanView    bool `json:"can_view"`
	CanAdd     bool `json:"can_add"`
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
	CanExport  bool `json:"can_export"`
}

type UpdateMenuAccessRequest struct {
	Items []UpdateMenuAccessItem `json:"items" validate:"required,dive"`
}
