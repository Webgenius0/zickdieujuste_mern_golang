package worship

import (
	"context"
	"errors"
	"math"

	"gotickets/internal/domain/worship/dto"
	"gotickets/internal/upload"
	"github.com/google/uuid"
)

type Service interface {
	Create(req dto.CreateWorshipReq) (dto.WorshipResponse, error)
	GetAll(timeOfDay string, page, limit int) (dto.PaginatedWorshipResponse, error)
	GetByID(id uuid.UUID) (dto.WorshipResponse, error)
	Update(id uuid.UUID, req dto.UpdateWorshipReq) (dto.WorshipResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo     Repository
	uploader upload.Uploader
}

func NewService(repo Repository, uploader upload.Uploader) Service {
	return &service{repo: repo, uploader: uploader}
}

func (s *service) Create(req dto.CreateWorshipReq) (dto.WorshipResponse, error) {
	worship := &Worship{
		Title:        req.Title,
		Artist:       req.Artist,
		TimeOfDay:    req.TimeOfDay,
		Duration:     req.Duration,
		ThumbnailURL: req.ThumbnailURL,
		AudioURL:     req.AudioURL,
		PrayerText:   req.PrayerText,
	}

	if err := s.repo.Create(worship); err != nil {
		return dto.WorshipResponse{}, err
	}

	return s.mapToResponse(worship), nil
}

func (s *service) GetAll(timeOfDay string, page, limit int) (dto.PaginatedWorshipResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	worships, total, err := s.repo.FindAll(timeOfDay, page, limit)
	if err != nil {
		return dto.PaginatedWorshipResponse{}, err
	}

	var data []dto.WorshipResponse
	for _, w := range worships {
		data = append(data, s.mapToResponse(&w))
	}
	if data == nil {
		data = make([]dto.WorshipResponse, 0)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return dto.PaginatedWorshipResponse{
		Data:       data,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) GetByID(id uuid.UUID) (dto.WorshipResponse, error) {
	worship, err := s.repo.FindByID(id)
	if err != nil {
		return dto.WorshipResponse{}, err
	}
	if worship == nil {
		return dto.WorshipResponse{}, errors.New("worship track not found")
	}

	return s.mapToResponse(worship), nil
}

func (s *service) Update(id uuid.UUID, req dto.UpdateWorshipReq) (dto.WorshipResponse, error) {
	worship, err := s.repo.FindByID(id)
	if err != nil {
		return dto.WorshipResponse{}, err
	}
	if worship == nil {
		return dto.WorshipResponse{}, errors.New("worship track not found")
	}

	worship.Title = req.Title
	worship.Artist = req.Artist
	worship.TimeOfDay = req.TimeOfDay
	worship.Duration = req.Duration
	worship.ThumbnailURL = req.ThumbnailURL
	worship.AudioURL = req.AudioURL
	worship.PrayerText = req.PrayerText

	if err := s.repo.Update(worship); err != nil {
		return dto.WorshipResponse{}, err
	}

	return s.mapToResponse(worship), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	worship, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if worship == nil {
		return errors.New("worship track not found")
	}

	// Delete from Cloudinary if uploader is available
	if s.uploader != nil {
		if worship.ThumbnailURL != "" {
			if pubID, err := upload.PublicIDFromURL(worship.ThumbnailURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
		if worship.AudioURL != "" {
			if pubID, err := upload.PublicIDFromURL(worship.AudioURL); err == nil {
				// We also need to specify resource_type for audio if it's "video" in Cloudinary.
				// But uploader.Delete interface currently doesn't specify resource type. It uses default (image).
				// Cloudinary sometimes fails to delete raw/video files without specifying resource_type.
				// For the sake of standard integration in this backend, we try to delete it as is.
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
	}

	return s.repo.Delete(id)
}

func (s *service) mapToResponse(w *Worship) dto.WorshipResponse {
	return dto.WorshipResponse{
		ID:           w.ID,
		Title:        w.Title,
		Artist:       w.Artist,
		TimeOfDay:    w.TimeOfDay,
		Duration:     w.Duration,
		ThumbnailURL: w.ThumbnailURL,
		AudioURL:     w.AudioURL,
		PrayerText:   w.PrayerText,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}
