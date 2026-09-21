package dto

import (
	"time"

	"github.com/google/uuid"
)

type LanguageResponse struct {
	ID        uuid.UUID `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	Name      string    `json:"name" example:"English"`
	Code      string    `json:"code" example:"en"`
	FlagIcon  string    `json:"flag_icon" example:"https://res.cloudinary.com/demo/image/upload/flag.png"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2026-08-17T15:00:00Z"`
}

type CreateLanguageReq struct {
	Name     string `json:"name" validate:"required" example:"English"`
	Code     string `json:"code" validate:"required" example:"en"`
	FlagIcon string `json:"flag_icon" example:"https://res.cloudinary.com/demo/image/upload/flag.png"`
	IsActive *bool  `json:"is_active" validate:"required" example:"true"`
}

type UpdateLanguageReq struct {
	Name     *string `json:"name" example:"English"`
	Code     *string `json:"code" example:"en"`
	FlagIcon *string `json:"flag_icon" example:"https://res.cloudinary.com/demo/image/upload/flag.png"`
	IsActive *bool   `json:"is_active" example:"true"`
}

type PaginatedLanguageResponse struct {
	Languages  []LanguageResponse `json:"languages"`
	TotalCount int64              `json:"total_count" example:"100"`
	Page       int                `json:"page" example:"1"`
	Limit      int                `json:"limit" example:"10"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
