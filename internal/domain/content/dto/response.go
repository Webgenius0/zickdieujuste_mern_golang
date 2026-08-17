package dto

import (
	"time"

	"github.com/google/uuid"
)

// ContentSummary is used in list responses.
type ContentSummary struct {
	ID              uuid.UUID  `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	Type            string     `json:"type" example:"PRAYER"`
	SubType         *string    `json:"sub_type" example:"Morning"`
	CategoryTag     *string    `json:"category_tag" example:"Thanksgiving"`
	Title           string     `json:"title" example:"Morning Prayer for Peace"`
	AuthorOrSpeaker *string    `json:"author_or_speaker" example:"Pastor John"`
	ThumbnailURL    *string    `json:"thumbnail_url" example:"https://example.com/thumb.jpg"`
	MediaType       string     `json:"media_type" example:"AUDIO"`
	DurationSeconds *int       `json:"duration_seconds" example:"180"`
	IsPremium       bool       `json:"is_premium" example:"true"`
	PublishedAt     *time.Time `json:"published_at" example:"2026-08-17T15:00:00Z"`
	Audiences       []string   `json:"audiences"`
}

// ContentDetail is used in single-item and daily-quote responses.
// MediaURL is omitted (set to nil) for premium content when the user is not premium.
type ContentDetail struct {
	ContentSummary
	BodyText *string          `json:"body_text" example:"Dear Lord, thank you for this day..."`
	MediaURL *string          `json:"media_url" example:"https://example.com/audio.mp3"` // nil if premium-gated
	Related  []ContentSummary `json:"related,omitempty"`
}

// PremiumGateResponse is returned (HTTP 403) when a non-premium user requests premium content.
type PremiumGateResponse struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"Operation successful"`
	Upgrade string `json:"upgrade_url" example:"/api/v1/subscriptions/plans"`
}

// ContentListResponse wraps a paginated list.
type ContentListResponse struct {
	Items    []ContentSummary `json:"items"`
	Total    int64            `json:"total" example:"100"`
	Page     int              `json:"page" example:"1"`
	PageSize int              `json:"page_size" example:"20"`
}
