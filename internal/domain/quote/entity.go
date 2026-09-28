package quote

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quote struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PublishDate time.Time      `gorm:"not null;index"`
	QuoteText   string         `gorm:"type:text;not null"`
	Reference   string         `gorm:"type:varchar(255)"`
	Explanation string         `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
