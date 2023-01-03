package request

import "time"

type TodoURI struct {
	ID uint64 `uri:"id" binding:"required"`
}

type TodoCreateRequest struct {
	ActivityGroupID uint64 `json:"activity_group_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
}

type TodoUpdateRequest struct {
	Title    string `json:"title,omitempty"`
	IsActive bool   `json:"is_active"`
}
type TodoResponse struct {
	ID         uint64     `json:"id"`
	Title      string     `json:"title"`
	ActivityID uint64     `json:"activity_group_id"`
	IsActive   string     `json:"is_active"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletetAt  *time.Time `json:"deleted_at"`
}

type TodoCreatedResponse struct {
	ID         uint64     `json:"id"`
	Title      string     `json:"title"`
	ActivityID uint64     `json:"activity_group_id"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletetAt  *time.Time `json:"deleted_at"`
}
