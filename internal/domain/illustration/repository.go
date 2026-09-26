package illustration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(item *Illustration) error
	FindAll(page, limit int) ([]Illustration, int64, error)
	FindByID(id uuid.UUID) (*Illustration, error)
	Update(item *Illustration) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(item *Illustration) error {
	return r.db.Create(item).Error
}

func (r *repository) FindAll(page, limit int) ([]Illustration, int64, error) {
	var items []Illustration
	var total int64

	query := r.db.Model(&Illustration{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Illustration, error) {
	var item Illustration
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) Update(item *Illustration) error {
	return r.db.Save(item).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Illustration{}, "id = ?", id).Error
}
