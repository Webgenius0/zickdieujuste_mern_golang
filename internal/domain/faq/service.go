package faq

import (
	"gotickets/internal/domain/faq/dto"
)

type Service interface {
	CreateFAQ(req dto.CreateFAQRequest) (*dto.FAQResponse, error)
	GetAdminFAQs() ([]dto.FAQResponse, error)
	GetPublicFAQs(search string) ([]dto.FAQResponse, error)
	UpdateFAQ(id string, req dto.UpdateFAQRequest) (*dto.FAQResponse, error)
	DeleteFAQ(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateFAQ(req dto.CreateFAQRequest) (*dto.FAQResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	faq := &FAQ{
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	}

	if err := s.repo.Create(faq); err != nil {
		return nil, err
	}

	return s.mapToResponse(faq), nil
}

func (s *service) GetAdminFAQs() ([]dto.FAQResponse, error) {
	faqs, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]dto.FAQResponse, 0)
	for _, f := range faqs {
		res = append(res, *s.mapToResponse(&f))
	}
	return res, nil
}

func (s *service) GetPublicFAQs(search string) ([]dto.FAQResponse, error) {
	faqs, err := s.repo.FindActiveWithSearch(search)
	if err != nil {
		return nil, err
	}

	res := make([]dto.FAQResponse, 0)
	for _, f := range faqs {
		res = append(res, *s.mapToResponse(&f))
	}
	return res, nil
}

func (s *service) UpdateFAQ(id string, req dto.UpdateFAQRequest) (*dto.FAQResponse, error) {
	faq, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Question != "" {
		faq.Question = req.Question
	}
	if req.Answer != "" {
		faq.Answer = req.Answer
	}
	if req.SortOrder != nil {
		faq.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		faq.IsActive = *req.IsActive
	}

	if err := s.repo.Update(faq); err != nil {
		return nil, err
	}

	return s.mapToResponse(faq), nil
}

func (s *service) DeleteFAQ(id string) error {
	return s.repo.Delete(id)
}

func (s *service) mapToResponse(faq *FAQ) *dto.FAQResponse {
	return &dto.FAQResponse{
		ID:        faq.ID.String(),
		Question:  faq.Question,
		Answer:    faq.Answer,
		SortOrder: faq.SortOrder,
		IsActive:  faq.IsActive,
		CreatedAt: faq.CreatedAt,
		UpdatedAt: faq.UpdatedAt,
	}
}
