package motivation

import (
	"errors"

	"gotickets/internal/domain/motivation/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrMotivationNotFound = errors.New("motivation not found")
)

type Service interface {
	Create(req dto.CreateMotivationReq) (*dto.MotivationResponse, error)
	FindAll() ([]*dto.MotivationResponse, error)
	FindByID(id uuid.UUID) (*dto.MotivationResponse, error)
	Update(id uuid.UUID, req dto.UpdateMotivationReq) (*dto.MotivationResponse, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req dto.CreateMotivationReq) (*dto.MotivationResponse, error) {
	m := &Motivation{
		Title:        req.Title,
		SpeakerName:  req.SpeakerName,
		Description:  req.Description,
		VideoURL:     req.VideoURL,
		ThumbnailURL: req.ThumbnailURL,
		Duration:     req.Duration,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}

	return toDTO(m), nil
}

func (s *service) FindAll() ([]*dto.MotivationResponse, error) {
	motivations, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]*dto.MotivationResponse, 0, len(motivations))
	for _, m := range motivations {
		res = append(res, toDTO(m))
	}
	return res, nil
}

func (s *service) FindByID(id uuid.UUID) (*dto.MotivationResponse, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMotivationNotFound
		}
		return nil, err
	}
	return toDTO(m), nil
}

func (s *service) Update(id uuid.UUID, req dto.UpdateMotivationReq) (*dto.MotivationResponse, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMotivationNotFound
		}
		return nil, err
	}

	m.Title = req.Title
	m.SpeakerName = req.SpeakerName
	m.Description = req.Description
	m.VideoURL = req.VideoURL
	m.ThumbnailURL = req.ThumbnailURL
	m.Duration = req.Duration

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}

	return toDTO(m), nil
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func toDTO(m *Motivation) *dto.MotivationResponse {
	return &dto.MotivationResponse{
		ID:           m.ID,
		Title:        m.Title,
		SpeakerName:  m.SpeakerName,
		Description:  m.Description,
		VideoURL:     m.VideoURL,
		ThumbnailURL: m.ThumbnailURL,
		Duration:     m.Duration,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
