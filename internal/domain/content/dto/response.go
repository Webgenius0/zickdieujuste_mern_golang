package dto

import (
	"time"

	"github.com/google/uuid"
)

// ContentSummary is used in list responses.
type ContentSummary struct {
	ID              uuid.UUID  `json:"id"`
	Type            string     `json:"type"`
	SubType         *string    `json:"sub_type"`
	CategoryTag     *string    `json:"category_tag"`
	Title           string     `json:"title"`
	AuthorOrSpeaker *string    `json:"author_or_speaker"`
	ThumbnailURL    *string    `json:"thumbnail_url"`
	MediaType       string     `json:"media_type"`
	DurationSeconds *int       `json:"duration_seconds"`
	IsPremium       bool       `json:"is_premium"`
	PublishedAt     *time.Time `json:"published_at"`
	Audiences       []string   `json:"audiences"`
}

// ContentDetail is used in single-item and daily-quote responses.
// MediaURL is omitted (set to nil) for premium content when the user is not premium.
type ContentDetail struct {
	ContentSummary
	BodyText   *string          `json:"body_text"`
	MediaURL   *string          `json:"media_url"` // nil if premium-gated
	Related    []ContentSummary `json:"related,omitempty"`
}

// PremiumGateResponse is returned (HTTP 403) when a non-premium user requests premium content.
type PremiumGateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Upgrade string `json:"upgrade_url"`
}

// ContentListResponse wraps a paginated list.
type ContentListResponse struct {
	Items    []ContentSummary `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}
