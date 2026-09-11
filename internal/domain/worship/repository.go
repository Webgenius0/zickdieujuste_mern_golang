package worship

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(worship *Worship) error
	FindAll(timeOfDay string, page, limit int) ([]Worship, int64, error)
	FindByID(id uuid.UUID) (*Worship, error)
	Update(worship *Worship) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(worship *Worship) error {
	return r.db.Create(worship).Error
}

func (r *repository) FindAll(timeOfDay string, page, limit int) ([]Worship, int64, error) {
	var worships []Worship
	var total int64

	query := r.db.Model(&Worship{})

	if timeOfDay != "" {
		query = query.Where("time_of_day = ?", timeOfDay)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&worships).Error; err != nil {
		return nil, 0, err
	}

	return worships, total, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Worship, error) {
	var worship Worship
	if err := r.db.First(&worship, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &worship, nil
}

func (r *repository) Update(worship *Worship) error {
	return r.db.Save(worship).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&Worship{}).Error
}
