package proverb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Proverb represents a proverb entry in the system.
type Proverb struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Title              string         `gorm:"type:varchar(255);not null"`
	Category           string         `gorm:"type:varchar(100);not null"`
	Duration           string         `gorm:"type:varchar(50);not null"`
	ThumbnailURL       string         `gorm:"type:text;not null"`
	AudioURL           string         `gorm:"type:text;not null"`
	ScriptureReference string         `gorm:"type:varchar(255);not null"`
	MainText           string         `gorm:"type:text;not null"`
	Explanation        string         `gorm:"type:text;not null"`
	PublishDate        time.Time      `gorm:"type:date;not null;index"`
	CreatedAt          time.Time      `gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate sets the UUID.
func (p *Proverb) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
