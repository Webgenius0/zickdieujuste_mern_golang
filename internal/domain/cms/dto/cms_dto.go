package dto

import "github.com/google/uuid"

type CMSPageSectionRequest struct {
	Heading   string `json:"heading" validate:"required"`
	Content   string `json:"content" validate:"required"`
	SortOrder int    `json:"sort_order"`
}

type CreateCMSPageRequest struct {
	Slug      string                  `json:"slug" validate:"required"`
	Title     string                  `json:"title" validate:"required"`
	IntroText string                  `json:"intro_text"`
	Sections  []CMSPageSectionRequest `json:"sections"`
}

type UpdateCMSPageRequest struct {
	Title     string                  `json:"title" validate:"required"`
	IntroText string                  `json:"intro_text"`
	Sections  []CMSPageSectionRequest `json:"sections"`
}

type CMSPageSectionResponse struct {
	ID        uuid.UUID `json:"id"`
	Heading   string    `json:"heading"`
	Content   string    `json:"content"`
	SortOrder int       `json:"sort_order"`
}

type CMSPageResponse struct {
	ID        uuid.UUID                `json:"id"`
	Slug      string                   `json:"slug"`
	Title     string                   `json:"title"`
	IntroText string                   `json:"intro_text"`
	Sections  []CMSPageSectionResponse `json:"sections"`
	CreatedAt string                   `json:"created_at"`
	UpdatedAt string                   `json:"updated_at"`
}
