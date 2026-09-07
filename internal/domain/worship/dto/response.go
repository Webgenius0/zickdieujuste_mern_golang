package dto

import (
	"time"

	"github.com/google/uuid"
)

type WorshipResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Artist       string    `json:"artist"`
	TimeOfDay    string    `json:"time_of_day"`
	Duration     string    `json:"duration"`
	ThumbnailURL string    `json:"thumbnail_url"`
	AudioURL     string    `json:"audio_url"`
	PrayerText   string    `json:"prayer_text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PaginatedWorshipResponse struct {
	Data       []WorshipResponse `json:"data"`
	TotalItems int64             `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}
