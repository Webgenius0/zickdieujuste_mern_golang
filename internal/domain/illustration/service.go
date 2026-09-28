package illustration

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gotickets/internal/upload"
)

type Service interface {
	Create(req CreateIllustrationReq) (IllustrationResponse, error)
	GetAll(page, limit int) (PaginatedIllustrationResponse, error)
	GetByID(id uuid.UUID) (IllustrationResponse, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateIllustrationReq) (IllustrationResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo     Repository
	uploader upload.Uploader
}

func NewService(repo Repository, uploader upload.Uploader) Service {
	return &service{repo: repo, uploader: uploader}
}

func (s *service) Create(req CreateIllustrationReq) (IllustrationResponse, error) {
	item := &Illustration{
		ContentText: req.ContentText,
		Reference:   req.Reference,
		AudioURL:    req.AudioURL,
	}

	if err := s.repo.Create(item); err != nil {
		return IllustrationResponse{}, err
	}

	return mapToResponse(item), nil
}

func (s *service) GetAll(page, limit int) (PaginatedIllustrationResponse, error) {
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
		return PaginatedIllustrationResponse{}, err
	}

	var data []IllustrationResponse
	for _, item := range items {
		data = append(data, mapToResponse(&item))
	}
	if data == nil {
		data = make([]IllustrationResponse, 0)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return PaginatedIllustrationResponse{
		Data:       data,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) GetByID(id uuid.UUID) (IllustrationResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return IllustrationResponse{}, err
	}
	if item == nil {
		return IllustrationResponse{}, errors.New("item not found")
	}
	return mapToResponse(item), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateIllustrationReq) (IllustrationResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return IllustrationResponse{}, err
	}
	if item == nil {
		return IllustrationResponse{}, errors.New("item not found")
	}

	// If audio URL is changing, delete the old audio from Cloudinary
	if s.uploader != nil && item.AudioURL != nil && req.AudioURL != nil && *item.AudioURL != *req.AudioURL && *item.AudioURL != "" {
		if pubID, err := upload.PublicIDFromURL(*item.AudioURL); err == nil {
			_ = s.uploader.Delete(ctx, pubID)
		}
	} else if s.uploader != nil && item.AudioURL != nil && req.AudioURL == nil && *item.AudioURL != "" {
		// Audio removed
		if pubID, err := upload.PublicIDFromURL(*item.AudioURL); err == nil {
			_ = s.uploader.Delete(ctx, pubID)
		}
	}

	item.ContentText = req.ContentText
	item.Reference = req.Reference
	item.AudioURL = req.AudioURL

	if err := s.repo.Update(item); err != nil {
		return IllustrationResponse{}, err
	}

	return mapToResponse(item), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}

	if s.uploader != nil && item.AudioURL != nil && *item.AudioURL != "" {
		if pubID, err := upload.PublicIDFromURL(*item.AudioURL); err == nil {
			_ = s.uploader.Delete(ctx, pubID)
		}
	}

	return s.repo.Delete(id)
}

func mapToResponse(item *Illustration) IllustrationResponse {
	return IllustrationResponse{
		ID:          item.ID,
		ContentText: item.ContentText,
		Reference:   item.Reference,
		AudioURL:    item.AudioURL,
		CreatedAt:   item.CreatedAt,
	}
}
