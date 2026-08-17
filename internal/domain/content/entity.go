package content

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContentType is the primary content category.
type ContentType string

// MediaType is the type of media attached to a content item.
type MediaType string

// Audience is the target audience for a content item.
type Audience string

const (
	ContentTypePrayer        ContentType = "PRAYER"
	ContentTypeMotivation    ContentType = "MOTIVATION"
	ContentTypeWorship       ContentType = "WORSHIP"
	ContentTypeProverb       ContentType = "PROVERB"
	ContentTypeDailyQuote    ContentType = "DAILY_QUOTE"
	ContentTypeIllustration  ContentType = "ILLUSTRATION"
	ContentTypeEncouragement ContentType = "ENCOURAGEMENT"

	MediaTypeAudio MediaType = "AUDIO"
	MediaTypeVideo MediaType = "VIDEO"
	MediaTypeNone  MediaType = "NONE"

	AudienceAll   Audience = "ALL"
	AudienceKids  Audience = "KIDS"
	AudienceTeens Audience = "TEENS"
)

// Content represents the CONTENT table (ERD §CONTENT).
type Content struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Type            ContentType `gorm:"type:varchar(50);not null;index"`
	SubType         *string     `gorm:"type:varchar(100)"` // Nullable: Morning/Night, Day/Night, etc.
	CategoryTag     *string     `gorm:"type:varchar(100)"` // Nullable: Thanksgiving, Intercession, etc.
	Title           string      `gorm:"type:varchar(500);not null"`
	AuthorOrSpeaker *string     `gorm:"type:varchar(255)"`
	BodyText        *string     `gorm:"type:text"` // scripture / prayer / article body
	MediaURL        *string     `gorm:"type:text"` // Cloudinary URL
	MediaType       MediaType   `gorm:"type:varchar(20);not null;default:'NONE'"`
	DurationSeconds *int
	ThumbnailURL    *string    `gorm:"type:text"`
	IsPremium       bool       `gorm:"not null;default:false;index"`
	PublishedAt     *time.Time `gorm:"index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relations
	Audiences      []ContentAudience `gorm:"foreignKey:ContentID"`
	RelatedContent []*Content        `gorm:"many2many:related_content;joinForeignKey:primary_content_id;joinReferences:related_content_id"`
}

// BeforeCreate assigns a UUID before inserting a new Content row.
func (c *Content) BeforeCreate(_ *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}

// ContentAudience represents the CONTENT_AUDIENCES join table (ERD §CONTENT_AUDIENCES).
type ContentAudience struct {
	ContentID uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false;index"`
	Audience  Audience  `gorm:"type:varchar(20);primaryKey;index"`
}

// RelatedContent is the self-referential join for RELATED_CONTENT (ERD §RELATED_CONTENT).
// GORM's many2many tag on Content handles this join table automatically.
// The struct below is defined explicitly for clarity and to allow AutoMigrate to see the table.
type RelatedContentJoin struct {
	PrimaryContentID uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false"`
	RelatedContentID uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false"`
}
