package faq

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(faq *FAQ) error
	FindAll() ([]FAQ, error)
	FindActiveWithSearch(search string) ([]FAQ, error)
	FindByID(id string) (*FAQ, error)
	Update(faq *FAQ) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(faq *FAQ) error {
	return r.db.Create(faq).Error
}

func (r *repository) FindAll() ([]FAQ, error) {
	var faqs []FAQ
	err := r.db.Order("sort_order ASC").Find(&faqs).Error
	return faqs, err
}

func (r *repository) FindActiveWithSearch(search string) ([]FAQ, error) {
	var faqs []FAQ
	query := r.db.Where("is_active = ?", true)

	if search != "" {
		// Use ILIKE for case-insensitive search in PostgreSQL
		query = query.Where("question ILIKE ? OR answer ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Order("sort_order ASC").Find(&faqs).Error
	return faqs, err
}

func (r *repository) FindByID(id string) (*FAQ, error) {
	var faq FAQ
	err := r.db.First(&faq, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &faq, nil
}

func (r *repository) Update(faq *FAQ) error {
	return r.db.Save(faq).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&FAQ{}, "id = ?", id).Error
}
