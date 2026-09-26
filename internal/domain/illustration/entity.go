package illustration

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Illustration struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ContentText string         `gorm:"type:text;not null"`
	Reference   string         `gorm:"type:varchar(255)"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
