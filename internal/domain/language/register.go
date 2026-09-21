package language

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
	adminGroup := e.Group("/api/v1/admin/languages", authMW, middlewares.RequireAdmin)
	adminGroup.POST("", h.CreateLanguage)
	adminGroup.GET("", h.GetAdminLanguages)
	adminGroup.PUT("/:id", h.UpdateLanguage)
	adminGroup.DELETE("/:id", h.DeleteLanguage)

	// Public Routes (Requires Authentication)
	publicGroup := e.Group("/api/v1/languages", authMW)
	publicGroup.GET("", h.GetPublicLanguages)
}
