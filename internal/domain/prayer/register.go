package prayer

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"gotickets/internal/auth"
	"gotickets/internal/middlewares"
	"gotickets/internal/upload"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService, uploader upload.Uploader) {
	repo := NewRepository(db)
	svc := NewService(repo, uploader)
	handler := NewHandler(svc)

	publicGroup := e.Group("/api/v1")

	// Public Categories
	publicGroup.GET("/categories", handler.GetAllCategories)
	publicGroup.GET("/categories/:id", handler.GetCategoryByID)

	// Public SubCategories
	publicGroup.GET("/subcategories", handler.GetSubCategoriesByCategoryID)
	publicGroup.GET("/subcategories/:id", handler.GetSubCategoryByID)

	// Public Prayers
	publicGroup.GET("/prayers", handler.GetAllPrayers)
	publicGroup.GET("/prayers/:id", handler.GetPrayerByID)
	publicGroup.GET("/prayers/mobile/metadata", handler.GetMobileMetadata)

	adminGroup := e.Group("/api/v1/admin")
	adminGroup.Use(middlewares.AuthMiddleware(jwtService))
	adminGroup.Use(middlewares.RequireAdmin)

	// Categories
	adminGroup.POST("/categories", handler.CreateCategory)
	adminGroup.GET("/categories", handler.GetAllCategories)
	adminGroup.GET("/categories/:id", handler.GetCategoryByID)
	adminGroup.PUT("/categories/:id", handler.UpdateCategory)
	adminGroup.DELETE("/categories/:id", handler.DeleteCategory)

	// SubCategories
	adminGroup.POST("/subcategories", handler.CreateSubCategory)
	adminGroup.GET("/subcategories", handler.GetSubCategoriesByCategoryID) // uses ?categoryId=
	adminGroup.GET("/subcategories/:id", handler.GetSubCategoryByID)
	adminGroup.PUT("/subcategories/:id", handler.UpdateSubCategory)
	adminGroup.DELETE("/subcategories/:id", handler.DeleteSubCategory)

	// Prayers
	adminGroup.POST("/prayers", handler.CreatePrayer)
	adminGroup.GET("/prayers", handler.GetAllPrayers)
	adminGroup.GET("/prayers/:id", handler.GetPrayerByID)
	adminGroup.PUT("/prayers/:id", handler.UpdatePrayer)
	adminGroup.DELETE("/prayers/:id", handler.DeletePrayer)
}
