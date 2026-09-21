package language

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Language represents the LANGUAGES table.
type Language struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Code      string         `gorm:"type:varchar(10);not null;uniqueIndex"`
	FlagIcon  string         `gorm:"type:varchar(255)"` // URL to flag image or country code
	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate assigns a UUID before inserting a new Language row.
func (l *Language) BeforeCreate(_ *gorm.DB) error {
	l.ID = uuid.New()
	return nil
}
