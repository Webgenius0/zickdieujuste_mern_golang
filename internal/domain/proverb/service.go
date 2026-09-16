package proverb

import (
	"context"
	"errors"
	"math"

	"gotickets/internal/domain/proverb/dto"
	"gotickets/internal/upload"

	"github.com/google/uuid"
)

type Service interface {
	Create(req dto.CreateProverbReq) (dto.ProverbResponse, error)
	GetAll(page, limit int, targetAudience string, excludeToday bool) (dto.PaginatedProverbResponse, error)
	GetToday(targetAudience string) (dto.MobileTodayResponse, error)
	GetByID(id uuid.UUID) (dto.ProverbResponse, error)
	Update(id uuid.UUID, req dto.UpdateProverbReq) (dto.ProverbResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo     Repository
	uploader upload.Uploader
}

func NewService(repo Repository, uploader upload.Uploader) Service {
	return &service{repo: repo, uploader: uploader}
}

func (s *service) Create(req dto.CreateProverbReq) (dto.ProverbResponse, error) {
	proverb := &Proverb{
		Title:              req.Title,
		Category:           req.Category,
		Duration:           req.Duration,
		ThumbnailURL:       req.ThumbnailURL,
		AudioURL:           req.AudioURL,
		ScriptureReference: req.ScriptureReference,
		MainText:           req.MainText,
		Explanation:        req.Explanation,
		TargetAudience:     TargetAudience(req.TargetAudience),
		PublishDate:        req.PublishDate,
	}

	if err := s.repo.Create(proverb); err != nil {
		return dto.ProverbResponse{}, err
	}

	return s.mapToResponse(proverb), nil
}

func (s *service) GetAll(page, limit int, targetAudience string, excludeToday bool) (dto.PaginatedProverbResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	proverbs, total, err := s.repo.FindAll(page, limit, targetAudience, excludeToday)
	if err != nil {
		return dto.PaginatedProverbResponse{}, err
	}

	var data []dto.ProverbResponse
	for _, p := range proverbs {
		data = append(data, s.mapToResponse(&p))
	}
	if data == nil {
		data = make([]dto.ProverbResponse, 0)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return dto.PaginatedProverbResponse{
		Data:       data,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) GetToday(targetAudience string) (dto.MobileTodayResponse, error) {
	proverb, err := s.repo.FindToday(targetAudience)
	if err != nil {
		return dto.MobileTodayResponse{}, err
	}
	if proverb == nil {
		return dto.MobileTodayResponse{}, errors.New("proverb not found")
	}

	todayResp := s.mapToResponse(proverb)

	// Fetch previous proverbs (page 1, limit 10, excluding today's proverb)
	prevProverbs, _, err := s.repo.FindAll(1, 10, targetAudience, true)
	if err != nil {
		return dto.MobileTodayResponse{}, err
	}

	var previous []dto.ProverbResponse
	for _, p := range prevProverbs {
		previous = append(previous, s.mapToResponse(&p))
	}
	if previous == nil {
		previous = make([]dto.ProverbResponse, 0)
	}

	return dto.MobileTodayResponse{
		Today:    todayResp,
		Previous: previous,
	}, nil
}

func (s *service) GetByID(id uuid.UUID) (dto.ProverbResponse, error) {
	proverb, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProverbResponse{}, err
	}
	if proverb == nil {
		return dto.ProverbResponse{}, errors.New("proverb not found")
	}

	return s.mapToResponse(proverb), nil
}

func (s *service) Update(id uuid.UUID, req dto.UpdateProverbReq) (dto.ProverbResponse, error) {
	proverb, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProverbResponse{}, err
	}
	if proverb == nil {
		return dto.ProverbResponse{}, errors.New("proverb not found")
	}

	proverb.Title = req.Title
	proverb.Category = req.Category
	proverb.Duration = req.Duration
	proverb.ThumbnailURL = req.ThumbnailURL
	proverb.AudioURL = req.AudioURL
	proverb.ScriptureReference = req.ScriptureReference
	proverb.MainText = req.MainText
	proverb.Explanation = req.Explanation
	proverb.TargetAudience = TargetAudience(req.TargetAudience)
	proverb.PublishDate = req.PublishDate

	if err := s.repo.Update(proverb); err != nil {
		return dto.ProverbResponse{}, err
	}

	return s.mapToResponse(proverb), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	proverb, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if proverb == nil {
		return errors.New("proverb not found")
	}

	// Delete from Cloudinary if uploader is available
	if s.uploader != nil {
		if proverb.ThumbnailURL != "" {
			if pubID, err := upload.PublicIDFromURL(proverb.ThumbnailURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
		if proverb.AudioURL != "" {
			if pubID, err := upload.PublicIDFromURL(proverb.AudioURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
	}

	return s.repo.Delete(id)
}

func (s *service) mapToResponse(p *Proverb) dto.ProverbResponse {
	return dto.ProverbResponse{
		ID:                 p.ID,
		Title:              p.Title,
		Category:           p.Category,
		Duration:           p.Duration,
		ThumbnailURL:       p.ThumbnailURL,
		AudioURL:           p.AudioURL,
		ScriptureReference: p.ScriptureReference,
		MainText:           p.MainText,
		Explanation:        p.Explanation,
		TargetAudience:     string(p.TargetAudience),
		PublishDate:        p.PublishDate,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}
