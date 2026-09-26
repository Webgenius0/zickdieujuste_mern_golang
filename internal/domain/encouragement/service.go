package encouragement

import (
	"errors"
	"math"

	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateEncouragementReq) (EncouragementResponse, error)
	GetAll(page, limit int) (PaginatedEncouragementResponse, error)
	GetByID(id uuid.UUID) (EncouragementResponse, error)
	Update(id uuid.UUID, req UpdateEncouragementReq) (EncouragementResponse, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req CreateEncouragementReq) (EncouragementResponse, error) {
	item := &Encouragement{
		ContentText: req.ContentText,
		Reference:   req.Reference,
	}

	if err := s.repo.Create(item); err != nil {
		return EncouragementResponse{}, err
	}

	return mapToResponse(item), nil
}

func (s *service) GetAll(page, limit int) (PaginatedEncouragementResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	items, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return PaginatedEncouragementResponse{}, err
	}

	var data []EncouragementResponse
	for _, item := range items {
		data = append(data, mapToResponse(&item))
	}
	if data == nil {
		data = make([]EncouragementResponse, 0)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return PaginatedEncouragementResponse{
		Data:       data,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) GetByID(id uuid.UUID) (EncouragementResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return EncouragementResponse{}, err
	}
	if item == nil {
		return EncouragementResponse{}, errors.New("item not found")
	}
	return mapToResponse(item), nil
}

func (s *service) Update(id uuid.UUID, req UpdateEncouragementReq) (EncouragementResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return EncouragementResponse{}, err
	}
	if item == nil {
		return EncouragementResponse{}, errors.New("item not found")
	}

	item.ContentText = req.ContentText
	item.Reference = req.Reference

	if err := s.repo.Update(item); err != nil {
		return EncouragementResponse{}, err
	}

	return mapToResponse(item), nil
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func mapToResponse(item *Encouragement) EncouragementResponse {
	return EncouragementResponse{
		ID:          item.ID,
		ContentText: item.ContentText,
		Reference:   item.Reference,
		CreatedAt:   item.CreatedAt,
	}
}
