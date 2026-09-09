package prayer

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gotickets/internal/domain/prayer/dto"
	"gotickets/internal/upload"
)

type Service interface {
	// Categories
	CreateCategory(req dto.CreateCategoryRequest) (dto.CategoryResponse, error)
	GetAllCategories() ([]dto.CategoryResponse, error)
	GetCategoryByID(id uuid.UUID) (dto.CategoryResponse, error)
	UpdateCategory(id uuid.UUID, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// SubCategories
	CreateSubCategory(req dto.CreateSubCategoryRequest) (dto.SubCategoryResponse, error)
	GetSubCategoriesByCategoryID(categoryID uuid.UUID) ([]dto.SubCategoryResponse, error)
	GetSubCategoryByID(id uuid.UUID) (dto.SubCategoryResponse, error)
	UpdateSubCategory(id uuid.UUID, req dto.UpdateSubCategoryRequest) (dto.SubCategoryResponse, error)
	DeleteSubCategory(ctx context.Context, id uuid.UUID) error

	// Prayers
	CreatePrayer(req dto.CreatePrayerRequest) (dto.PrayerResponse, error)
	GetAllPrayers(page, limit int, filters map[string]interface{}) (map[string]interface{}, error)
	GetPrayerByID(id uuid.UUID) (dto.PrayerResponse, error)
	UpdatePrayer(ctx context.Context, id uuid.UUID, req dto.UpdatePrayerRequest) (dto.PrayerResponse, error)
	DeletePrayer(ctx context.Context, id uuid.UUID) error

	// Mobile
	GetMobileMetadata(targetAudience string) (dto.MobileMetadataResponse, error)
}

type service struct {
	repo     Repository
	uploader upload.Uploader
}

func NewService(repo Repository, uploader upload.Uploader) Service {
	return &service{repo: repo, uploader: uploader}
}

// Helpers
func mapCategoryToResponse(c *Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:             c.ID,
		Name:           c.Name,
		TargetAudience: string(c.TargetAudience),
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func mapSubCategoryToResponse(s *SubCategory) dto.SubCategoryResponse {
	return dto.SubCategoryResponse{
		ID:         s.ID,
		CategoryID: s.CategoryID,
		Name:       s.Name,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func mapPrayerToResponse(p *Prayer) dto.PrayerResponse {
	resp := dto.PrayerResponse{
		ID:            p.ID,
		Title:         p.Title,
		CategoryID:    p.CategoryID,
		SubCategoryID: p.SubCategoryID,
		AgeGroup:      p.AgeGroup,
		MediaType:     string(p.MediaType),
		Duration:      p.Duration,
		ThumbnailURL:  p.ThumbnailURL,
		MediaURL:      p.MediaURL,
		ContentText:   p.ContentText,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
	
	if p.Category != nil {
		catResp := mapCategoryToResponse(p.Category)
		resp.Category = &catResp
	}
	if p.SubCategory != nil {
		subCatResp := mapSubCategoryToResponse(p.SubCategory)
		resp.SubCategory = &subCatResp
	}
	
	return resp
}

// Categories
func (s *service) CreateCategory(req dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	cat := &Category{
		Name:           req.Name,
		TargetAudience: TargetAudience(req.TargetAudience),
	}
	if err := s.repo.CreateCategory(cat); err != nil {
		return dto.CategoryResponse{}, err
	}
	return mapCategoryToResponse(cat), nil
}

func (s *service) GetAllCategories() ([]dto.CategoryResponse, error) {
	cats, err := s.repo.FindAllCategories()
	if err != nil {
		return nil, err
	}
	var res []dto.CategoryResponse
	for _, c := range cats {
		res = append(res, mapCategoryToResponse(&c))
	}
	if res == nil {
		res = make([]dto.CategoryResponse, 0)
	}
	return res, nil
}

func (s *service) GetCategoryByID(id uuid.UUID) (dto.CategoryResponse, error) {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	if cat == nil {
		return dto.CategoryResponse{}, errors.New("category not found")
	}
	return mapCategoryToResponse(cat), nil
}

func (s *service) UpdateCategory(id uuid.UUID, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	if cat == nil {
		return dto.CategoryResponse{}, errors.New("category not found")
	}

	cat.Name = req.Name
	cat.TargetAudience = TargetAudience(req.TargetAudience)

	if err := s.repo.UpdateCategory(cat); err != nil {
		return dto.CategoryResponse{}, err
	}
	return mapCategoryToResponse(cat), nil
}

func (s *service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategory(id)
}

// SubCategories
func (s *service) CreateSubCategory(req dto.CreateSubCategoryRequest) (dto.SubCategoryResponse, error) {
	sub := &SubCategory{
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}
	if err := s.repo.CreateSubCategory(sub); err != nil {
		return dto.SubCategoryResponse{}, err
	}
	return mapSubCategoryToResponse(sub), nil
}

func (s *service) GetSubCategoriesByCategoryID(categoryID uuid.UUID) ([]dto.SubCategoryResponse, error) {
	subs, err := s.repo.FindSubCategoriesByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}
	var res []dto.SubCategoryResponse
	for _, sub := range subs {
		res = append(res, mapSubCategoryToResponse(&sub))
	}
	if res == nil {
		res = make([]dto.SubCategoryResponse, 0)
	}
	return res, nil
}

func (s *service) GetSubCategoryByID(id uuid.UUID) (dto.SubCategoryResponse, error) {
	sub, err := s.repo.FindSubCategoryByID(id)
	if err != nil {
		return dto.SubCategoryResponse{}, err
	}
	if sub == nil {
		return dto.SubCategoryResponse{}, errors.New("subcategory not found")
	}
	return mapSubCategoryToResponse(sub), nil
}

func (s *service) UpdateSubCategory(id uuid.UUID, req dto.UpdateSubCategoryRequest) (dto.SubCategoryResponse, error) {
	sub, err := s.repo.FindSubCategoryByID(id)
	if err != nil {
		return dto.SubCategoryResponse{}, err
	}
	if sub == nil {
		return dto.SubCategoryResponse{}, errors.New("subcategory not found")
	}

	sub.CategoryID = req.CategoryID
	sub.Name = req.Name

	if err := s.repo.UpdateSubCategory(sub); err != nil {
		return dto.SubCategoryResponse{}, err
	}
	return mapSubCategoryToResponse(sub), nil
}

func (s *service) DeleteSubCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSubCategory(id)
}

// Prayers
func (s *service) CreatePrayer(req dto.CreatePrayerRequest) (dto.PrayerResponse, error) {
	p := &Prayer{
		Title:         req.Title,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		AgeGroup:      req.AgeGroup,
		MediaType:     MediaType(req.MediaType),
		Duration:      req.Duration,
		ThumbnailURL:  req.ThumbnailURL,
		MediaURL:      req.MediaURL,
		ContentText:   req.ContentText,
	}
	if err := s.repo.CreatePrayer(p); err != nil {
		return dto.PrayerResponse{}, err
	}
	
	loadedP, _ := s.repo.FindPrayerByID(p.ID)
	if loadedP != nil {
		return mapPrayerToResponse(loadedP), nil
	}
	return mapPrayerToResponse(p), nil
}

func (s *service) GetAllPrayers(page, limit int, filters map[string]interface{}) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	prayers, total, err := s.repo.FindAllPrayers(page, limit, filters)
	if err != nil {
		return nil, err
	}

	var data []dto.PrayerResponse
	for _, p := range prayers {
		data = append(data, mapPrayerToResponse(&p))
	}
	if data == nil {
		data = make([]dto.PrayerResponse, 0)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return map[string]interface{}{
		"data":       data,
		"totalItems": total,
		"totalPages": totalPages,
		"page":       page,
		"limit":      limit,
	}, nil
}

func (s *service) GetPrayerByID(id uuid.UUID) (dto.PrayerResponse, error) {
	p, err := s.repo.FindPrayerByID(id)
	if err != nil {
		return dto.PrayerResponse{}, err
	}
	if p == nil {
		return dto.PrayerResponse{}, errors.New("prayer not found")
	}
	return mapPrayerToResponse(p), nil
}

func (s *service) UpdatePrayer(ctx context.Context, id uuid.UUID, req dto.UpdatePrayerRequest) (dto.PrayerResponse, error) {
	p, err := s.repo.FindPrayerByID(id)
	if err != nil {
		return dto.PrayerResponse{}, err
	}
	if p == nil {
		return dto.PrayerResponse{}, errors.New("prayer not found")
	}

	if s.uploader != nil {
		if p.ThumbnailURL != "" && p.ThumbnailURL != req.ThumbnailURL {
			if pubID, err := upload.PublicIDFromURL(p.ThumbnailURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
		if p.MediaURL != "" && p.MediaURL != req.MediaURL {
			if pubID, err := upload.PublicIDFromURL(p.MediaURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
	}

	p.Title = req.Title
	p.CategoryID = req.CategoryID
	p.SubCategoryID = req.SubCategoryID
	p.AgeGroup = req.AgeGroup
	p.MediaType = MediaType(req.MediaType)
	p.Duration = req.Duration
	p.ThumbnailURL = req.ThumbnailURL
	p.MediaURL = req.MediaURL
	p.ContentText = req.ContentText

	if err := s.repo.UpdatePrayer(p); err != nil {
		return dto.PrayerResponse{}, err
	}
	
	loadedP, _ := s.repo.FindPrayerByID(p.ID)
	if loadedP != nil {
		return mapPrayerToResponse(loadedP), nil
	}
	
	return mapPrayerToResponse(p), nil
}

func (s *service) DeletePrayer(ctx context.Context, id uuid.UUID) error {
	p, err := s.repo.FindPrayerByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("prayer not found")
	}

	if s.uploader != nil {
		if p.ThumbnailURL != "" {
			if pubID, err := upload.PublicIDFromURL(p.ThumbnailURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
		if p.MediaURL != "" {
			if pubID, err := upload.PublicIDFromURL(p.MediaURL); err == nil {
				_ = s.uploader.Delete(ctx, pubID)
			}
		}
	}

	return s.repo.DeletePrayer(id)
}

// Mobile
func (s *service) GetMobileMetadata(targetAudience string) (dto.MobileMetadataResponse, error) {
	var cats []Category
	if err := s.repo.(*repository).db.Where("target_audience = ?", targetAudience).Order("name ASC").Find(&cats).Error; err != nil {
		return dto.MobileMetadataResponse{}, err
	}

	var mobileCats []dto.MobileCategoryResponse
	for _, c := range cats {
		var subs []SubCategory
		if err := s.repo.(*repository).db.Where("category_id = ?", c.ID).Order("name ASC").Find(&subs).Error; err != nil {
			return dto.MobileMetadataResponse{}, err
		}
		var subResponses []dto.SubCategoryResponse
		for _, sub := range subs {
			subResponses = append(subResponses, mapSubCategoryToResponse(&sub))
		}
		if subResponses == nil {
			subResponses = make([]dto.SubCategoryResponse, 0)
		}

		mobileCats = append(mobileCats, dto.MobileCategoryResponse{
			ID:             c.ID,
			Name:           c.Name,
			TargetAudience: string(c.TargetAudience),
			SubCategories:  subResponses,
		})
	}
	if mobileCats == nil {
		mobileCats = make([]dto.MobileCategoryResponse, 0)
	}

	// Predefine age groups for the mobile UI
	ageGroups := []string{"Age 0-5", "Age 6-13", "Age 13-15", "Age 15-18"}

	return dto.MobileMetadataResponse{
		TargetAudience: targetAudience,
		Categories:     mobileCats,
		AgeGroups:      ageGroups,
	}, nil
}
