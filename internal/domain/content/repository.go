package content

import (
	"errors"
	"time"

	"gotickets/internal/domain/content/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the data access contract for the content domain.
type Repository interface {
	List(filter dto.ListContentRequest) ([]Content, int64, error)
	GetByID(id uuid.UUID) (*Content, error)
	GetDailyQuote(date time.Time) (*Content, error)
	GetRelated(id uuid.UUID) ([]Content, error)
}

// GORM implementation

type repository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed content Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(filter dto.ListContentRequest) ([]Content, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	q := r.db.Model(&Content{}).Preload("Audiences")

	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.SubType != "" {
		q = q.Where("sub_type = ?", filter.SubType)
	}
	if filter.CategoryTag != "" {
		q = q.Where("category_tag = ?", filter.CategoryTag)
	}
	if filter.Audience != "" {
		// Join through content_audiences to filter by audience
		q = q.Joins("JOIN content_audiences ca ON ca.content_id = contents.id").
			Where("ca.audience = ?", filter.Audience)
	}
	if filter.Q != "" {
		// Full-text search via the generated search_vector column.
		// Falls back gracefully to ILIKE if the tsvector migration hasn't been applied yet.
		q = q.Where(
			"search_vector @@ plainto_tsquery('english', ?) OR title ILIKE ?",
			filter.Q, "%"+filter.Q+"%",
		)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []Content
	if err := q.Order("published_at DESC").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) GetByID(id uuid.UUID) (*Content, error) {
	var c Content
	err := r.db.Preload("Audiences").
		Preload("RelatedContent").
		Where("id = ?", id).
		First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("content not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetDailyQuote(date time.Time) (*Content, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var c Content
	err := r.db.Preload("Audiences").
		Where("type = ? AND published_at >= ? AND published_at < ?", ContentTypeDailyQuote, start, end).
		Order("published_at DESC").
		First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no daily quote for today")
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetRelated(id uuid.UUID) ([]Content, error) {
	var primary Content
	err := r.db.Preload("RelatedContent.Audiences").
		Where("id = ?", id).
		First(&primary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("content not found")
		}
		return nil, err
	}
	// Dereference related content pointers
	related := make([]Content, 0, len(primary.RelatedContent))
	for _, rc := range primary.RelatedContent {
		if rc != nil {
			related = append(related, *rc)
		}
	}
	return related, nil
}
