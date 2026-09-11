package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubCategoryRequest struct {
	CategoryID uuid.UUID `json:"categoryId" validate:"required"`
	Name       string    `json:"name" validate:"required"`
}

type UpdateSubCategoryRequest struct {
	CategoryID uuid.UUID `json:"categoryId" validate:"required"`
	Name       string    `json:"name" validate:"required"`
}

type SubCategoryResponse struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"categoryId"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
