package dto

import (
	"time"

	"github.com/google/uuid"
)

// PlanResponse is the shape of a single subscription plan.
type PlanResponse struct {
	ID              uuid.UUID `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	PriceAmount     float64   `json:"price_amount"`
	Currency        string    `json:"currency"`
	BillingInterval string    `json:"billing_interval"`
}

// SubscriptionResponse is the shape returned after verifying a receipt.
type SubscriptionResponse struct {
	ID                    uuid.UUID    `json:"id"`
	UserID                uuid.UUID    `json:"user_id"`
	Plan                  PlanResponse `json:"plan"`
	Store                 string       `json:"store"`
	Status                string       `json:"status"`
	ExternalTransactionID string       `json:"external_transaction_id"`
	StartDate             time.Time    `json:"start_date"`
	ExpiresAt             time.Time    `json:"expires_at"`
}

// MessageResponse is a generic success message.
type MessageResponse struct {
	Message string `json:"message"`
}
