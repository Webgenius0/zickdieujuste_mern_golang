package worship

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Worship represents a worship audio track in the system.
type Worship struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title        string         `gorm:"type:varchar(255);not null"`
	Artist       string         `gorm:"type:varchar(255);not null"`
	TimeOfDay    string         `gorm:"type:varchar(50);not null"` // e.g., "Day", "Night"
	Duration     string         `gorm:"type:varchar(50);not null"` // e.g., "5:42"
	ThumbnailURL string         `gorm:"type:text;not null"`
	AudioURL     string         `gorm:"type:text;not null"`
	PrayerText   string         `gorm:"type:text"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
