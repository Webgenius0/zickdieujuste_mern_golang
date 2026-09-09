package prayer

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Categories
	CreateCategory(category *Category) error
	FindAllCategories() ([]Category, error)
	FindCategoryByID(id uuid.UUID) (*Category, error)
	UpdateCategory(category *Category) error
	DeleteCategory(id uuid.UUID) error

	// SubCategories
	CreateSubCategory(sub *SubCategory) error
	FindSubCategoriesByCategoryID(categoryID uuid.UUID) ([]SubCategory, error)
	FindSubCategoryByID(id uuid.UUID) (*SubCategory, error)
	UpdateSubCategory(sub *SubCategory) error
	DeleteSubCategory(id uuid.UUID) error

	// Prayers
	CreatePrayer(prayer *Prayer) error
	FindAllPrayers(page, limit int, filters map[string]interface{}) ([]Prayer, int64, error)
	FindPrayerByID(id uuid.UUID) (*Prayer, error)
	UpdatePrayer(prayer *Prayer) error
	DeletePrayer(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Categories
func (r *repository) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}

func (r *repository) FindAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) FindCategoryByID(id uuid.UUID) (*Category, error) {
	var category Category
	if err := r.db.First(&category, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *repository) UpdateCategory(category *Category) error {
	return r.db.Save(category).Error
}

func (r *repository) DeleteCategory(id uuid.UUID) error {
	return r.db.Delete(&Category{}, "id = ?", id).Error
}

// SubCategories
func (r *repository) CreateSubCategory(sub *SubCategory) error {
	return r.db.Create(sub).Error
}

func (r *repository) FindSubCategoriesByCategoryID(categoryID uuid.UUID) ([]SubCategory, error) {
	var subs []SubCategory
	if err := r.db.Where("category_id = ?", categoryID).Order("name ASC").Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *repository) FindSubCategoryByID(id uuid.UUID) (*SubCategory, error) {
	var sub SubCategory
	if err := r.db.First(&sub, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *repository) UpdateSubCategory(sub *SubCategory) error {
	return r.db.Save(sub).Error
}

func (r *repository) DeleteSubCategory(id uuid.UUID) error {
	return r.db.Delete(&SubCategory{}, "id = ?", id).Error
}

// Prayers
func (r *repository) CreatePrayer(prayer *Prayer) error {
	return r.db.Create(prayer).Error
}

func (r *repository) FindAllPrayers(page, limit int, filters map[string]interface{}) ([]Prayer, int64, error) {
	var prayers []Prayer
	var total int64

	query := r.db.Model(&Prayer{})

	// Apply joins if needed for filtering
	if targetAudience, ok := filters["targetAudience"]; ok && targetAudience != "" {
		query = query.Joins("JOIN categories ON categories.id = prayers.category_id").
			Where("categories.target_audience = ?", targetAudience)
	}

	if categoryID, ok := filters["categoryId"]; ok && categoryID != "" {
		query = query.Where("prayers.category_id = ?", categoryID)
	}
	if ageGroup, ok := filters["ageGroup"]; ok && ageGroup != "" {
		query = query.Where("prayers.age_group = ?", ageGroup)
	}
	if mediaType, ok := filters["mediaType"]; ok && mediaType != "" {
		query = query.Where("prayers.media_type = ?", mediaType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Category").Preload("SubCategory").
		Order("prayers.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&prayers).Error; err != nil {
		return nil, 0, err
	}

	return prayers, total, nil
}

func (r *repository) FindPrayerByID(id uuid.UUID) (*Prayer, error) {
	var p Prayer
	if err := r.db.Preload("Category").Preload("SubCategory").First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *repository) UpdatePrayer(prayer *Prayer) error {
	return r.db.Save(prayer).Error
}

func (r *repository) DeletePrayer(id uuid.UUID) error {
	return r.db.Delete(&Prayer{}, "id = ?", id).Error
}
