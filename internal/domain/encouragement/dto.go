package encouragement

import (
	"time"
	"github.com/google/uuid"
)

type CreateEncouragementReq struct {
	ContentText string `json:"contentText" validate:"required"`
	Reference   string `json:"reference"`
}

type UpdateEncouragementReq struct {
	ContentText string `json:"contentText" validate:"required"`
	Reference   string `json:"reference"`
}

type EncouragementResponse struct {
	ID          uuid.UUID `json:"id"`
	ContentText string    `json:"contentText"`
	Reference   string    `json:"reference"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PaginatedEncouragementResponse struct {
	Data       []EncouragementResponse `json:"data"`
	TotalItems int64                 `json:"totalItems"`
	TotalPages int                   `json:"totalPages"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
}
