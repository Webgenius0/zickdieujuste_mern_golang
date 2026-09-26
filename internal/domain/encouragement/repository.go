package encouragement

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(item *Encouragement) error
	FindAll(page, limit int) ([]Encouragement, int64, error)
	FindByID(id uuid.UUID) (*Encouragement, error)
	Update(item *Encouragement) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(item *Encouragement) error {
	return r.db.Create(item).Error
}

func (r *repository) FindAll(page, limit int) ([]Encouragement, int64, error) {
	var items []Encouragement
	var total int64

	query := r.db.Model(&Encouragement{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Encouragement, error) {
	var item Encouragement
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) Update(item *Encouragement) error {
	return r.db.Save(item).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Encouragement{}, "id = ?", id).Error
}
