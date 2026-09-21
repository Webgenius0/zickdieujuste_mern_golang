package language

import (
	"errors"
	"fmt"
	"gotickets/internal/domain/language/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrLanguageNotFound = errors.New("language not found")
	ErrLanguageConflict = errors.New("language with this code already exists")
)

type Service interface {
	CreateLanguage(req dto.CreateLanguageReq) (*dto.LanguageResponse, error)
	GetLanguages(page, limit int, search string, activeOnly bool) (*dto.PaginatedLanguageResponse, error)
	GetLanguageByID(id uuid.UUID) (*dto.LanguageResponse, error)
	UpdateLanguage(id uuid.UUID, req dto.UpdateLanguageReq) (*dto.LanguageResponse, error)
	DeleteLanguage(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateLanguage(req dto.CreateLanguageReq) (*dto.LanguageResponse, error) {
	existing, err := s.repo.FindByCode(req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing language: %w", err)
	}
	if existing != nil {
		return nil, ErrLanguageConflict
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	lang := &Language{
		Name:     req.Name,
		Code:     req.Code,
		FlagIcon: req.FlagIcon,
		IsActive: isActive,
	}

	if err := s.repo.Create(lang); err != nil {
		return nil, fmt.Errorf("failed to create language: %w", err)
	}

	return toLanguageResponse(lang), nil
}

func (s *service) GetLanguages(page, limit int, search string, activeOnly bool) (*dto.PaginatedLanguageResponse, error) {
	languages, total, err := s.repo.FindAll(page, limit, search, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch languages: %w", err)
	}

	var dtos []dto.LanguageResponse
	for _, lang := range languages {
		dtos = append(dtos, *toLanguageResponse(&lang))
	}

	if dtos == nil {
		dtos = []dto.LanguageResponse{}
	}

	return &dto.PaginatedLanguageResponse{
		Languages:  dtos,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) GetLanguageByID(id uuid.UUID) (*dto.LanguageResponse, error) {
	lang, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLanguageNotFound
		}
		return nil, fmt.Errorf("failed to fetch language: %w", err)
	}

	return toLanguageResponse(lang), nil
}

func (s *service) UpdateLanguage(id uuid.UUID, req dto.UpdateLanguageReq) (*dto.LanguageResponse, error) {
	lang, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLanguageNotFound
		}
		return nil, fmt.Errorf("failed to fetch language: %w", err)
	}

	if req.Code != nil && *req.Code != lang.Code {
		existing, err := s.repo.FindByCode(*req.Code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check existing language: %w", err)
		}
		if existing != nil && existing.ID != lang.ID {
			return nil, ErrLanguageConflict
		}
		lang.Code = *req.Code
	}

	if req.Name != nil {
		lang.Name = *req.Name
	}
	if req.FlagIcon != nil {
		lang.FlagIcon = *req.FlagIcon
	}
	if req.IsActive != nil {
		lang.IsActive = *req.IsActive
	}

	if err := s.repo.Update(lang); err != nil {
		return nil, fmt.Errorf("failed to update language: %w", err)
	}

	return toLanguageResponse(lang), nil
}

func (s *service) DeleteLanguage(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLanguageNotFound
		}
		return fmt.Errorf("failed to fetch language: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete language: %w", err)
	}

	return nil
}

func toLanguageResponse(lang *Language) *dto.LanguageResponse {
	return &dto.LanguageResponse{
		ID:        lang.ID,
		Name:      lang.Name,
		Code:      lang.Code,
		FlagIcon:  lang.FlagIcon,
		IsActive:  lang.IsActive,
		CreatedAt: lang.CreatedAt,
	}
}
