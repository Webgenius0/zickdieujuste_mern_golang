package dto

import (
	"time"

	"github.com/google/uuid"
)

type MotivationResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	SpeakerName  string    `json:"speaker_name"`
	Description  string    `json:"description"`
	VideoURL     string    `json:"video_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Duration     string    `json:"duration"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
