package illustration

import (
	"time"
	"github.com/google/uuid"
)

type CreateIllustrationReq struct {
	ContentText string  `json:"contentText" validate:"required"`
	Reference   string  `json:"reference"`
	AudioURL    *string `json:"audioUrl"`
}

type UpdateIllustrationReq struct {
	ContentText string  `json:"contentText" validate:"required"`
	Reference   string  `json:"reference"`
	AudioURL    *string `json:"audioUrl"`
}

type IllustrationResponse struct {
	ID          uuid.UUID `json:"id"`
	ContentText string    `json:"contentText"`
	Reference   string    `json:"reference"`
	AudioURL    *string   `json:"audioUrl"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PaginatedIllustrationResponse struct {
	Data       []IllustrationResponse `json:"data"`
	TotalItems int64                 `json:"totalItems"`
	TotalPages int                   `json:"totalPages"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
}
