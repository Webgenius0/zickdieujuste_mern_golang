package motivation

import (
	"gotickets/internal/querybuilder"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(m *Motivation) error
	FindAll(params querybuilder.Params) ([]*Motivation, *querybuilder.Meta, error)
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

func (r *repository) FindAll(params querybuilder.Params) ([]*Motivation, *querybuilder.Meta, error) {
	var motivations []*Motivation

	meta, err := querybuilder.New(r.db.Model(&Motivation{}), params).
		Search([]string{"title", "speaker_name"}).
		Filter().
		Sort().
		Paginate().
		ExecuteWithMeta(&motivations)

	if err != nil {
		return nil, nil, err
	}
	return motivations, meta, nil
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
