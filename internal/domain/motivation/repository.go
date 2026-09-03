package motivation

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(m *Motivation) error
	FindAll() ([]*Motivation, error)
	FindByID(id uuid.UUID) (*Motivation, error)
	Update(m *Motivation) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(m *Motivation) error {
	return r.db.Create(m).Error
}

func (r *repository) FindAll() ([]*Motivation, error) {
	var motivations []*Motivation
	err := r.db.Order("created_at desc").Find(&motivations).Error
	return motivations, err
}

func (r *repository) FindByID(id uuid.UUID) (*Motivation, error) {
	var m Motivation
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) Update(m *Motivation) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&Motivation{}).Error
}
