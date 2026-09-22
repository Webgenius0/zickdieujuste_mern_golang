package cms

import (
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, authMW echo.MiddlewareFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	// Admin Routes (Requires Authentication and Admin Role)
	adminGroup := e.Group("/api/v1/admin/cms/pages", authMW, middlewares.RequireAdmin)
	adminGroup.POST("", h.CreateAdminPage)
	adminGroup.GET("", h.GetAdminPages)
	adminGroup.GET("/:slug", h.GetAdminPageBySlug)
	adminGroup.PUT("/:slug", h.UpdateAdminPage)
	adminGroup.DELETE("/:slug", h.DeleteAdminPage)

	// Public Routes (Optional Authentication or No Authentication depending on requirements)
	// For CMS pages like Privacy Policy, they are typically public
	publicGroup := e.Group("/api/v1/cms/pages")
	publicGroup.GET("/privacy-policy", h.GetPrivacyPolicy)
	publicGroup.GET("/terms-and-conditions", h.GetTermsAndConditions)
}
