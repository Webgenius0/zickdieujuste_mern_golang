package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePrayerRequest struct {
	Title         string     `json:"title" validate:"required"`
	CategoryID    uuid.UUID  `json:"categoryId" validate:"required"`
	SubCategoryID *uuid.UUID `json:"subCategoryId"`
	AgeGroup      *string    `json:"ageGroup"`
	MediaType     string     `json:"mediaType" validate:"required"`
	Duration      string     `json:"duration"`
	ThumbnailURL  string     `json:"thumbnailUrl"`
	MediaURL      string     `json:"mediaUrl"`
	ContentText   string     `json:"contentText"`
}

type UpdatePrayerRequest struct {
	Title         string     `json:"title" validate:"required"`
	CategoryID    uuid.UUID  `json:"categoryId" validate:"required"`
	SubCategoryID *uuid.UUID `json:"subCategoryId"`
	AgeGroup      *string    `json:"ageGroup"`
	MediaType     string     `json:"mediaType" validate:"required"`
	Duration      string     `json:"duration"`
	ThumbnailURL  string     `json:"thumbnailUrl"`
	MediaURL      string     `json:"mediaUrl"`
	ContentText   string     `json:"contentText"`
}

type PrayerResponse struct {
	ID            uuid.UUID            `json:"id"`
	Title         string               `json:"title"`
	CategoryID    uuid.UUID            `json:"categoryId"`
	SubCategoryID *uuid.UUID           `json:"subCategoryId"`
	AgeGroup      *string              `json:"ageGroup"`
	MediaType     string               `json:"mediaType"`
	Duration      string               `json:"duration"`
	ThumbnailURL  string               `json:"thumbnailUrl"`
	MediaURL      string               `json:"mediaUrl"`
	ContentText   string               `json:"contentText"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	
	// Preloaded relations
	Category      *CategoryResponse    `json:"category,omitempty"`
	SubCategory   *SubCategoryResponse `json:"subCategory,omitempty"`
}

type PaginatedPrayerResponse struct {
	Data       []PrayerResponse `json:"data"`
	TotalItems int              `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
}
