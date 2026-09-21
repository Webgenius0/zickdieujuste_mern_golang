package faq

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
	adminGroup := e.Group("/api/v1/admin/faqs", authMW, middlewares.RequireAdmin)
	adminGroup.POST("", h.CreateFAQ)
	adminGroup.GET("", h.GetAdminFAQs)
	adminGroup.PUT("/:id", h.UpdateFAQ)
	adminGroup.DELETE("/:id", h.DeleteFAQ)

	// Public Routes (Requires Authentication)
	publicGroup := e.Group("/api/v1/faqs", authMW)
	publicGroup.GET("", h.GetPublicFAQs)
}
