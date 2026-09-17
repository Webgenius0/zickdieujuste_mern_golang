package proverb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(proverb *Proverb) error
	FindAll(page, limit int, targetAudience string, excludeToday bool) ([]Proverb, int64, error)
	FindToday(targetAudience string) (*Proverb, error)
	FindByID(id uuid.UUID) (*Proverb, error)
	Update(proverb *Proverb) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(proverb *Proverb) error {
	return r.db.Create(proverb).Error
}

func (r *repository) FindAll(page, limit int, targetAudience string, excludeToday bool) ([]Proverb, int64, error) {
	var proverbs []Proverb
	var total int64

	query := r.db.Model(&Proverb{})
	if targetAudience != "" {
		query = query.Where("target_audience = ?", targetAudience)
	}

	if excludeToday {
		var todayProverb Proverb
		todayQuery := r.db.Model(&Proverb{})
		if targetAudience != "" {
			todayQuery = todayQuery.Where("target_audience = ?", targetAudience)
		}
		// If we find a "today" proverb, exclude its ID from the main query
		if err := todayQuery.Order("publish_date DESC").First(&todayProverb).Error; err == nil {
			query = query.Where("id != ?", todayProverb.ID)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("publish_date DESC").Offset(offset).Limit(limit).Find(&proverbs).Error; err != nil {
		return nil, 0, err
	}

	return proverbs, total, nil
}

func (r *repository) FindToday(targetAudience string) (*Proverb, error) {
	var proverb Proverb
	query := r.db.Model(&Proverb{})
	if targetAudience != "" {
		query = query.Where("target_audience = ?", targetAudience)
	}
	if err := query.Order("publish_date DESC").First(&proverb).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &proverb, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Proverb, error) {
	var proverb Proverb
	if err := r.db.Where("id = ?", id).First(&proverb).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &proverb, nil
}

func (r *repository) Update(proverb *Proverb) error {
	return r.db.Save(proverb).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Proverb{}, id).Error
}
