package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name           string `json:"name" validate:"required"`
	TargetAudience string `json:"targetAudience" validate:"required"`
	Module         string `json:"module" validate:"omitempty,oneof='Prayer' 'Faith'"`
}

type UpdateCategoryRequest struct {
	Name           string `json:"name" validate:"required"`
	TargetAudience string `json:"targetAudience" validate:"required"`
	Module         string `json:"module" validate:"omitempty,oneof='Prayer' 'Faith'"`
}

type CategoryResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	TargetAudience string    `json:"targetAudience"`
	Module         string    `json:"module"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
