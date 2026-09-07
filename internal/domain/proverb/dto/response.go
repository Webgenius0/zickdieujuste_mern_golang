package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProverbResponse struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `json:"title"`
	Category           string    `json:"category"`
	Duration           string    `json:"duration"`
	ThumbnailURL       string    `json:"thumbnail_url"`
	AudioURL           string    `json:"audio_url"`
	ScriptureReference string    `json:"scripture_reference"`
	MainText           string    `json:"main_text"`
	Explanation        string    `json:"explanation"`
	PublishDate        time.Time `json:"publish_date"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PaginatedProverbResponse struct {
	Data       []ProverbResponse `json:"data"`
	TotalItems int64             `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}
