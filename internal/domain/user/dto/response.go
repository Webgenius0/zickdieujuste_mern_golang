package dto

import (
	"time"

	"github.com/google/uuid"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
}

type ProfileResponse struct {
	ID                 uuid.UUID  `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	Name               string     `json:"name" example:"John Doe"`
	Email              string     `json:"email" example:"user@example.com"`
	AuthProvider       string     `json:"auth_provider" example:"EMAIL"`
	Location           *string    `json:"location" example:"New York, USA"`
	AvatarURL          *string    `json:"avatar_url" example:"https://res.cloudinary.com/demo/image/upload/avatar.jpg"`
	ThemePreference    string     `json:"theme_preference" example:"NAVY"`
	LanguagePreference string     `json:"language_preference" example:"en"`
	Age                int        `json:"age" example:"25"`
	IsPremium          bool       `json:"is_premium" example:"true"`
	TermsAcceptedAt    *time.Time `json:"terms_accepted_at" example:"2026-08-17T15:00:00Z"`
	CreatedAt          time.Time  `json:"created_at" example:"2026-08-17T15:00:00Z"`
}

type AvatarResponse struct {
	AvatarURL string `json:"avatar_url" example:"https://res.cloudinary.com/demo/image/upload/avatar.jpg"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
