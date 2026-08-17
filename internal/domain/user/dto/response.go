package dto

import (
	"time"

	"github.com/google/uuid"
)

// AuthResponse is returned by register, login, and refresh endpoints.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ProfileResponse is returned by GET /api/v1/users/me and PUT /api/v1/users/me.
type ProfileResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Email              string     `json:"email"`
	AuthProvider       string     `json:"auth_provider"`
	Location           *string    `json:"location"`
	AvatarURL          *string    `json:"avatar_url"`
	ThemePreference    string     `json:"theme_preference"`
	LanguagePreference string     `json:"language_preference"`
	IsPremium          bool       `json:"is_premium"`
	TermsAcceptedAt    *time.Time `json:"terms_accepted_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

// AvatarResponse is returned by POST /api/v1/users/me/avatar.
type AvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
}

// MessageResponse is a generic success message response.
type MessageResponse struct {
	Message string `json:"message"`
}
