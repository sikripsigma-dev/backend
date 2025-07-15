package dto

type ChatRooms struct {
	ID          string             `json:"id"`
	StudentID   string             `json:"student_id"`
	CompanyID   string             `json:"company_id"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	Company     *CompanySimple     `json:"company,omitempty"`
	LastMessage *ChatMessageSimple `json:"last_message,omitempty"`
}

type CompanySimple struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChatMessageSimple struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type ChatRoomsWithLatestMessage struct {
	ID            string `json:"id"`
	StudentID     string `json:"student_id"`
	CompanyID     string `json:"company_id"`
	CompanyName   string `json:"company_name"`
	LastMessage   string `json:"last_message,omitempty"`
	LastMessageAt string `json:"last_message_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	UnreadCount   int    `json:"unread_count"`
}

type ChatMessageResponse struct {
	ID         string  `json:"id"`
	RoomID     string  `json:"room_id"`
	SenderID   *string `json:"sender_id"`
	SenderType string  `json:"sender_type"`
	Message    string  `json:"message"`
	IsRead     bool    `json:"is_read"`
	CreatedAt  string  `json:"created_at"`
}
