package subscription

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the data access contract for the subscription domain.
type Repository interface {
	GetActivePlans() ([]SubscriptionPlan, error)
	GetPlanByCode(code PlanCode) (*SubscriptionPlan, error)
	UpsertSubscription(s *Subscription) error
	GetActiveSubscription(userID uuid.UUID) (*Subscription, error)
	CountPlans() (int64, error)
	CreatePlan(p *SubscriptionPlan) error
}

// ──────────────────────────────────────────────────────────────────────────────
// GORM implementation
// ──────────────────────────────────────────────────────────────────────────────

type repository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed subscription Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetActivePlans() ([]SubscriptionPlan, error) {
	var plans []SubscriptionPlan
	if err := r.db.Where("is_active = true").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *repository) GetPlanByCode(code PlanCode) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := r.db.Where("code = ?", code).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan not found")
		}
		return nil, err
	}
	return &plan, nil
}

func (r *repository) UpsertSubscription(s *Subscription) error {
	return r.db.Save(s).Error
}

func (r *repository) GetActiveSubscription(userID uuid.UUID) (*Subscription, error) {
	var sub Subscription
	err := r.db.Preload("Plan").
		Where("user_id = ? AND status = ?", userID, StatusActive).
		Order("expires_at DESC").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *repository) CountPlans() (int64, error) {
	var count int64
	return count, r.db.Model(&SubscriptionPlan{}).Count(&count).Error
}

func (r *repository) CreatePlan(p *SubscriptionPlan) error {
	return r.db.Create(p).Error
}
