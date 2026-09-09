package prayer

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TargetAudience string

const (
	AudienceGeneral TargetAudience = "General"
	AudienceKids    TargetAudience = "Kids"
	AudienceTeens   TargetAudience = "Teens"
)

type MediaType string

const (
	MediaTypeAudio MediaType = "Audio"
	MediaTypeVideo MediaType = "Video"
)

type Category struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name           string         `gorm:"type:varchar(100);not null"`
	TargetAudience TargetAudience `gorm:"type:varchar(50);not null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	SubCategories []SubCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

type SubCategory struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name       string         `gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Category *Category `gorm:"foreignKey:CategoryID"`
}

func (s *SubCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type Prayer struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Title         string         `gorm:"type:varchar(255);not null"`
	CategoryID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	SubCategoryID *uuid.UUID     `gorm:"type:uuid;index"`
	AgeGroup      *string        `gorm:"type:varchar(50)"`
	MediaType     MediaType      `gorm:"type:varchar(20);not null"`
	Duration      string         `gorm:"type:varchar(50)"`
	ThumbnailURL  string         `gorm:"type:text"`
	MediaURL      string         `gorm:"type:text"`
	ContentText   string         `gorm:"type:text"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Category    *Category    `gorm:"foreignKey:CategoryID"`
	SubCategory *SubCategory `gorm:"foreignKey:SubCategoryID"`
}

func (p *Prayer) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
