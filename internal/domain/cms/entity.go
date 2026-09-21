package cms

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSPage struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Slug      string           `gorm:"uniqueIndex;not null" json:"slug"`
	Title     string           `gorm:"not null" json:"title"`
	IntroText string           `gorm:"type:text" json:"intro_text"`
	Sections  []CMSPageSection `gorm:"foreignKey:PageID;constraint:OnDelete:CASCADE" json:"sections"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}

type CMSPageSection struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PageID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"page_id"`
	Heading   string         `json:"heading"`
	Content   string         `gorm:"type:text" json:"content"`
	SortOrder int            `json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
