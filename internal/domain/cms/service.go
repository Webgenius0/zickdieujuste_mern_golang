package cms

import (
	"gotickets/internal/domain/cms/dto"

	"github.com/google/uuid"
)

type Service interface {
	GetAllPages() ([]dto.CMSPageResponse, error)
	GetPageBySlug(slug string) (*dto.CMSPageResponse, error)
	UpdatePage(slug string, req dto.UpdateCMSPageRequest) (*dto.CMSPageResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllPages() ([]dto.CMSPageResponse, error) {
	pages, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]dto.CMSPageResponse, 0)
	for _, p := range pages {
		res = append(res, *s.mapToResponse(&p))
	}
	return res, nil
}

func (s *service) GetPageBySlug(slug string) (*dto.CMSPageResponse, error) {
	page, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(page), nil
}

func (s *service) UpdatePage(slug string, req dto.UpdateCMSPageRequest) (*dto.CMSPageResponse, error) {
	page, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	page.Title = req.Title
	page.IntroText = req.IntroText

	sections := make([]CMSPageSection, 0)
	for _, reqSec := range req.Sections {
		sections = append(sections, CMSPageSection{
			ID:        uuid.New(),
			PageID:    page.ID,
			Heading:   reqSec.Heading,
			Content:   reqSec.Content,
			SortOrder: reqSec.SortOrder,
		})
	}

	if err := s.repo.UpdatePageWithSections(page, sections); err != nil {
		return nil, err
	}

	// Refetch to get the updated state with correct timestamps
	updatedPage, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(updatedPage), nil
}

func (s *service) mapToResponse(page *CMSPage) *dto.CMSPageResponse {
	sectionRes := make([]dto.CMSPageSectionResponse, 0)
	for _, sec := range page.Sections {
		sectionRes = append(sectionRes, dto.CMSPageSectionResponse{
			ID:        sec.ID,
			Heading:   sec.Heading,
			Content:   sec.Content,
			SortOrder: sec.SortOrder,
		})
	}

	return &dto.CMSPageResponse{
		ID:        page.ID,
		Slug:      page.Slug,
		Title:     page.Title,
		IntroText: page.IntroText,
		Sections:  sectionRes,
		CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
