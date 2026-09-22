package quote

import (
	"github.com/google/uuid"
)

type QuoteResponse struct {
	ID          uuid.UUID `json:"id"`
	PublishDate string    `json:"publish_date"`
	QuoteText   string    `json:"quote_text"`
	Reference   string    `json:"reference"`
	Explanation string    `json:"explanation"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type CreateQuoteRequest struct {
	PublishDate string `json:"publish_date" validate:"required"` // Format YYYY-MM-DD
	QuoteText   string `json:"quote_text" validate:"required"`
	Reference   string `json:"reference"`
	Explanation string `json:"explanation"`
}

type UpdateQuoteRequest struct {
	PublishDate string `json:"publish_date" validate:"required"`
	QuoteText   string `json:"quote_text" validate:"required"`
	Reference   string `json:"reference"`
	Explanation string `json:"explanation"`
}
