package content

import (
	"errors"
	"time"

	"gotickets/internal/domain/content/dto"

	"github.com/google/uuid"
)

// Service defines the business logic for the content domain.
type Service interface {
	List(filter dto.ListContentRequest) (*dto.ContentListResponse, error)
	GetByID(id uuid.UUID, isPremiumUser bool) (*dto.ContentDetail, error)
	GetDailyQuote(isPremiumUser bool) (*dto.ContentDetail, error)
	GetRelated(id uuid.UUID) ([]dto.ContentSummary, error)
}

type service struct {
	repo Repository
}

// NewService creates a new content Service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(filter dto.ListContentRequest) (*dto.ContentListResponse, error) {
	items, total, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	summaries := make([]dto.ContentSummary, 0, len(items))
	for _, c := range items {
		summaries = append(summaries, toSummary(c))
	}

	return &dto.ContentListResponse{
		Items:    summaries,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *service) GetByID(id uuid.UUID, isPremiumUser bool) (*dto.ContentDetail, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toDetail(*c, isPremiumUser), nil
}

func (s *service) GetDailyQuote(isPremiumUser bool) (*dto.ContentDetail, error) {
	c, err := s.repo.GetDailyQuote(time.Now())
	if err != nil {
		return nil, errors.New("no daily quote available for today")
	}
	return toDetail(*c, isPremiumUser), nil
}

func (s *service) GetRelated(id uuid.UUID) ([]dto.ContentSummary, error) {
	items, err := s.repo.GetRelated(id)
	if err != nil {
		return nil, err
	}
	summaries := make([]dto.ContentSummary, 0, len(items))
	for _, c := range items {
		summaries = append(summaries, toSummary(c))
	}
	return summaries, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Mapper helpers
// ──────────────────────────────────────────────────────────────────────────────

func toSummary(c Content) dto.ContentSummary {
	audiences := make([]string, 0, len(c.Audiences))
	for _, a := range c.Audiences {
		audiences = append(audiences, string(a.Audience))
	}
	return dto.ContentSummary{
		ID:              c.ID,
		Type:            string(c.Type),
		SubType:         c.SubType,
		CategoryTag:     c.CategoryTag,
		Title:           c.Title,
		AuthorOrSpeaker: c.AuthorOrSpeaker,
		ThumbnailURL:    c.ThumbnailURL,
		MediaType:       string(c.MediaType),
		DurationSeconds: c.DurationSeconds,
		IsPremium:       c.IsPremium,
		PublishedAt:     c.PublishedAt,
		Audiences:       audiences,
	}
}

func toDetail(c Content, isPremiumUser bool) *dto.ContentDetail {
	d := &dto.ContentDetail{
		ContentSummary: toSummary(c),
		BodyText:       c.BodyText,
	}

	// Premium gating: only expose media_url if user is premium or content is free
	if !c.IsPremium || isPremiumUser {
		d.MediaURL = c.MediaURL
	}

	return d
}
