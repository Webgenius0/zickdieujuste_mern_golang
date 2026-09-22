package quote

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Service interface {
	CreateQuote(req CreateQuoteRequest) (*QuoteResponse, error)
	GetAllQuotes() ([]QuoteResponse, error)
	GetPublishedQuotes() ([]QuoteResponse, error)
	GetPublishedQuoteByID(id uuid.UUID) (*QuoteResponse, error)
	UpdateQuote(id uuid.UUID, req UpdateQuoteRequest) (*QuoteResponse, error)
	DeleteQuote(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateQuote(req CreateQuoteRequest) (*QuoteResponse, error) {
	publishDate, err := time.Parse("2006-01-02", req.PublishDate)
	if err != nil {
		return nil, err
	}

	quote := &Quote{
		PublishDate: publishDate,
		QuoteText:   req.QuoteText,
		Reference:   req.Reference,
		Explanation: req.Explanation,
	}

	if err := s.repo.Create(quote); err != nil {
		return nil, err
	}

	return s.mapToResponse(quote), nil
}

func (s *service) GetAllQuotes() ([]QuoteResponse, error) {
	quotes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]QuoteResponse, 0, len(quotes))
	for _, q := range quotes {
		res = append(res, *s.mapToResponse(&q))
	}
	return res, nil
}

func (s *service) GetPublishedQuotes() ([]QuoteResponse, error) {
	quotes, err := s.repo.FindPublished()
	if err != nil {
		return nil, err
	}

	res := make([]QuoteResponse, 0, len(quotes))
	for _, q := range quotes {
		res = append(res, *s.mapToResponse(&q))
	}
	return res, nil
}

func (s *service) GetPublishedQuoteByID(id uuid.UUID) (*QuoteResponse, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify it is published (PublishDate <= Now)
	if quote.PublishDate.After(time.Now()) {
		return nil, echo.ErrNotFound // Or a custom not published error
	}

	return s.mapToResponse(quote), nil
}

func (s *service) UpdateQuote(id uuid.UUID, req UpdateQuoteRequest) (*QuoteResponse, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	publishDate, err := time.Parse("2006-01-02", req.PublishDate)
	if err != nil {
		return nil, err
	}

	quote.PublishDate = publishDate
	quote.QuoteText = req.QuoteText
	quote.Reference = req.Reference
	quote.Explanation = req.Explanation

	if err := s.repo.Update(quote); err != nil {
		return nil, err
	}

	return s.mapToResponse(quote), nil
}

func (s *service) DeleteQuote(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) mapToResponse(quote *Quote) *QuoteResponse {
	return &QuoteResponse{
		ID:          quote.ID,
		PublishDate: quote.PublishDate.Format("2006-01-02"),
		QuoteText:   quote.QuoteText,
		Reference:   quote.Reference,
		Explanation: quote.Explanation,
		CreatedAt:   quote.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   quote.UpdatedAt.Format(time.RFC3339),
	}
}
