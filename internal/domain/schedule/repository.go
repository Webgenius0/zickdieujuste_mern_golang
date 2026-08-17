package schedule

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the data access contract for the schedule domain.
type Repository interface {
	GetByUserID(userID uuid.UUID) (*UserSchedule, error)
	Upsert(s *UserSchedule) error
}

// ──────────────────────────────────────────────────────────────────────────────
// GORM implementation
// ──────────────────────────────────────────────────────────────────────────────

type repository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed schedule Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByUserID(userID uuid.UUID) (*UserSchedule, error) {
	var s UserSchedule
	err := r.db.Where("user_id = ?", userID).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // caller creates default on nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) Upsert(s *UserSchedule) error {
	// Save handles both insert (new) and update (existing) via primary key presence
	return r.db.Save(s).Error
}
