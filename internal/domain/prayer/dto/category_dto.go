package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name           string `json:"name" validate:"required"`
	TargetAudience string `json:"targetAudience" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name           string `json:"name" validate:"required"`
	TargetAudience string `json:"targetAudience" validate:"required"`
}

type CategoryResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	TargetAudience string    `json:"targetAudience"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
