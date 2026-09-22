package quote

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(quote *Quote) error
	FindAll() ([]Quote, error)
	FindPublished() ([]Quote, error)
	FindByID(id uuid.UUID) (*Quote, error)
	Update(quote *Quote) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(quote *Quote) error {
	return r.db.Create(quote).Error
}

func (r *repository) FindAll() ([]Quote, error) {
	var quotes []Quote
	err := r.db.Order("publish_date DESC").Find(&quotes).Error
	return quotes, err
}

func (r *repository) FindPublished() ([]Quote, error) {
	var quotes []Quote
	err := r.db.Where("publish_date <= ?", time.Now()).Order("publish_date DESC").Find(&quotes).Error
	return quotes, err
}

func (r *repository) FindByID(id uuid.UUID) (*Quote, error) {
	var quote Quote
	err := r.db.First(&quote, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *repository) Update(quote *Quote) error {
	return r.db.Save(quote).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Quote{}, "id = ?", id).Error
}
