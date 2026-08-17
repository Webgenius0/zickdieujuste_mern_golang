package dto

import (
	"time"

	"github.com/google/uuid"
)

// PlanResponse is the shape of a single subscription plan.
type PlanResponse struct {
	ID              uuid.UUID `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	Code            string    `json:"code" example:"PREMIUM_REQUIRED"`
	Name            string    `json:"name" example:"John Doe"`
	PriceAmount     float64   `json:"price_amount" example:"9.99"`
	Currency        string    `json:"currency" example:"USD"`
	BillingInterval string    `json:"billing_interval" example:"MONTHLY"`
}

// SubscriptionResponse is the shape returned after verifying a receipt.
type SubscriptionResponse struct {
	ID                    uuid.UUID    `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	UserID                uuid.UUID    `json:"user_id" example:"b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e"`
	Plan                  PlanResponse `json:"plan"`
	Store                 string       `json:"store" example:"APPLE"`
	Status                string       `json:"status" example:"ACTIVE"`
	ExternalTransactionID string       `json:"external_transaction_id" example:"1000000123456"`
	StartDate             time.Time    `json:"start_date" example:"2026-08-17T15:00:00Z"`
	ExpiresAt             time.Time    `json:"expires_at" example:"2026-09-17T15:00:00Z"`
}

// MessageResponse is a generic success message.
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
