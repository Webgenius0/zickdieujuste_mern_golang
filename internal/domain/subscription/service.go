package subscription

import (
	"fmt"
	"time"

	"gotickets/internal/domain/subscription/dto"
	"gotickets/internal/domain/user"

	"github.com/google/uuid"
)

// ReceiptVerifier is the interface for verifying store receipts.
// Actual Apple/Google validation is out of scope for this pass — the no-op stub is wired in.
// To add real validation: implement this interface (e.g., in internal/platform/apple/ or internal/platform/google/)
// and swap the wiring in register.go — no service or handler code changes required.
type ReceiptVerifier interface {
	Verify(store Store, receiptPayload string) (*VerifiedReceipt, error)
}

// VerifiedReceipt is the normalized result returned by a ReceiptVerifier.
type VerifiedReceipt struct {
	TransactionID string
	PlanCode      PlanCode
	ExpiresAt     time.Time
}

// noopVerifier is the stub implementation wired in until real store integration is available.
type noopVerifier struct{}

func (n *noopVerifier) Verify(_ Store, _ string) (*VerifiedReceipt, error) {
	// TODO: replace with real Apple/Google receipt validation
	return &VerifiedReceipt{
		TransactionID: "stub-txn-" + uuid.New().String(),
		PlanCode:      PlanCodeAnnual,
		ExpiresAt:     time.Now().Add(365 * 24 * time.Hour),
	}, nil
}

// NewNoopVerifier returns the stub ReceiptVerifier.
func NewNoopVerifier() ReceiptVerifier {
	return &noopVerifier{}
}

// ──────────────────────────────────────────────────────────────────────────────
// Service
// ──────────────────────────────────────────────────────────────────────────────

// Service defines the business logic for the subscription domain.
type Service interface {
	ListPlans() ([]dto.PlanResponse, error)
	VerifyReceipt(userID uuid.UUID, req dto.VerifyReceiptRequest) (*dto.SubscriptionResponse, error)
	HandleWebhook(req dto.WebhookRequest) error
}

type service struct {
	repo     Repository
	verifier ReceiptVerifier
	userRepo user.Repository
}

// NewService creates a new subscription Service.
// userRepo is required to update users.is_premium after successful verification.
func NewService(repo Repository, verifier ReceiptVerifier, userRepo user.Repository) Service {
	return &service{repo: repo, verifier: verifier, userRepo: userRepo}
}

func (s *service) ListPlans() ([]dto.PlanResponse, error) {
	plans, err := s.repo.GetActivePlans()
	if err != nil {
		return nil, err
	}
	resp := make([]dto.PlanResponse, 0, len(plans))
	for _, p := range plans {
		resp = append(resp, toPlanResponse(p))
	}
	return resp, nil
}

func (s *service) VerifyReceipt(userID uuid.UUID, req dto.VerifyReceiptRequest) (*dto.SubscriptionResponse, error) {
	verified, err := s.verifier.Verify(Store(req.Store), req.ReceiptPayload)
	if err != nil {
		return nil, fmt.Errorf("receipt verification failed: %w", err)
	}

	plan, err := s.repo.GetPlanByCode(verified.PlanCode)
	if err != nil {
		return nil, fmt.Errorf("plan not found for code %q: %w", verified.PlanCode, err)
	}

	sub := &Subscription{
		UserID:                userID,
		PlanID:                plan.ID,
		Store:                 Store(req.Store),
		Status:                StatusActive,
		ExternalTransactionID: verified.TransactionID,
		StartDate:             time.Now(),
		ExpiresAt:             verified.ExpiresAt,
		Plan:                  *plan,
	}
	if err := s.repo.UpsertSubscription(sub); err != nil {
		return nil, fmt.Errorf("failed to persist subscription: %w", err)
	}

	// Mark user as premium
	u, err := s.userRepo.GetUserByID(userID)
	if err == nil {
		u.IsPremium = true
		_ = s.userRepo.UpdateUser(u)
	}

	return toSubResponse(*sub), nil
}

func (s *service) HandleWebhook(req dto.WebhookRequest) error {
	// TODO: verify store signature before processing (Apple/Google webhook signing key setup)
	// Processing logic for renewal/cancel/refund events goes here
	// For now, log and acknowledge receipt
	fmt.Printf("[Webhook] store=%s payload_keys=%d\n", req.Store, len(req.Payload))
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Seed: run once at startup if subscription_plans table is empty
// ──────────────────────────────────────────────────────────────────────────────

// SeedPlans inserts the 3 known ZICK subscription plans if the table is empty.
// Called from register.go after AutoMigrate.
func SeedPlans(repo Repository) error {
	count, err := repo.CountPlans()
	if err != nil || count > 0 {
		return err // already seeded or error
	}

	plans := []SubscriptionPlan{
		{Code: PlanCodeBiannual, Name: "Biannual", PriceAmount: 4.99, Currency: "USD", BillingInterval: IntervalMonth},
		{Code: PlanCodeAnnual, Name: "Annual", PriceAmount: 39.99, Currency: "USD", BillingInterval: IntervalYear},
		{Code: PlanCodeFamily, Name: "Friends & Family", PriceAmount: 79.99, Currency: "USD", BillingInterval: IntervalMonth},
	}
	for i := range plans {
		if err := repo.CreatePlan(&plans[i]); err != nil {
			return fmt.Errorf("failed to seed plan %q: %w", plans[i].Code, err)
		}
	}
	fmt.Println("[Seed] Subscription plans seeded successfully")
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Mapper helpers
// ──────────────────────────────────────────────────────────────────────────────

func toPlanResponse(p SubscriptionPlan) dto.PlanResponse {
	return dto.PlanResponse{
		ID:              p.ID,
		Code:            string(p.Code),
		Name:            p.Name,
		PriceAmount:     p.PriceAmount,
		Currency:        p.Currency,
		BillingInterval: string(p.BillingInterval),
	}
}

func toSubResponse(s Subscription) *dto.SubscriptionResponse {
	return &dto.SubscriptionResponse{
		ID:                    s.ID,
		UserID:                s.UserID,
		Plan:                  toPlanResponse(s.Plan),
		Store:                 string(s.Store),
		Status:                string(s.Status),
		ExternalTransactionID: s.ExternalTransactionID,
		StartDate:             s.StartDate,
		ExpiresAt:             s.ExpiresAt,
	}
}
