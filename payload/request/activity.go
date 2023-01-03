package request

import "time"

type ActivityIdURI struct {
	ID uint64 `uri:"id" binding:"required"`
}
type ActivityRequest struct {
	Title string `json:"title" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type ActivityUpdateRequest struct {
	Title string `json:"title" binding:"required"`
}

type ActivityCreateResponse struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ActivityGetOneResponse struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
