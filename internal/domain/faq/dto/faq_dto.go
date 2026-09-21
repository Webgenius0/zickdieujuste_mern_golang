package dto

import "time"

type CreateFAQRequest struct {
	Question  string `json:"question" validate:"required"`
	Answer    string `json:"answer" validate:"required"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"` // Pointer to allow false
}

type UpdateFAQRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder *int   `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

type FAQResponse struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
