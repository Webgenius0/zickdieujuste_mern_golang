package language

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(lang *Language) error
	FindAll(page, limit int, search string, activeOnly bool) ([]Language, int64, error)
	FindByID(id uuid.UUID) (*Language, error)
	FindByCode(code string) (*Language, error)
	Update(lang *Language) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(lang *Language) error {
	return r.db.Create(lang).Error
}

func (r *repository) FindAll(page, limit int, search string, activeOnly bool) ([]Language, int64, error) {
	var languages []Language
	var total int64

	query := r.db.Model(&Language{})

	if search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&languages).Error; err != nil {
		return nil, 0, err
	}

	return languages, total, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Language, error) {
	var lang Language
	if err := r.db.First(&lang, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lang, nil
}

func (r *repository) FindByCode(code string) (*Language, error) {
	var lang Language
	if err := r.db.First(&lang, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &lang, nil
}

func (r *repository) Update(lang *Language) error {
	return r.db.Save(lang).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Language{}, "id = ?", id).Error
}
