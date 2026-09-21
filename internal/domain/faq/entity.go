package faq

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FAQ represents the Help & Support questions and answers.
type FAQ struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Question   string         `gorm:"type:varchar(255);not null"`
	Answer     string         `gorm:"type:text;not null"`
	SortOrder  int            `gorm:"not null;default:0"`
	IsActive   bool           `gorm:"not null;default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate assigns a UUID before inserting a new FAQ row.
func (f *FAQ) BeforeCreate(_ *gorm.DB) error {
	f.ID = uuid.New()
	return nil
}
