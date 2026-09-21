package cms

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]CMSPage, error)
	FindBySlug(slug string) (*CMSPage, error)
	UpdatePageWithSections(page *CMSPage, sections []CMSPageSection) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]CMSPage, error) {
	var pages []CMSPage
	err := r.db.Order("slug ASC").Find(&pages).Error
	return pages, err
}

func (r *repository) FindBySlug(slug string) (*CMSPage, error) {
	var page CMSPage
	err := r.db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&page, "slug = ?", slug).Error
	
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *repository) UpdatePageWithSections(page *CMSPage, sections []CMSPageSection) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update the page details (Title, IntroText)
		if err := tx.Save(page).Error; err != nil {
			return err
		}

		// Delete existing sections for this page
		if err := tx.Where("page_id = ?", page.ID).Delete(&CMSPageSection{}).Error; err != nil {
			return err
		}

		// Insert new sections
		if len(sections) > 0 {
			if err := tx.Create(&sections).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
