package quote

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
	adminGroup := e.Group("/api/v1/admin/quotes", authMW, middlewares.RequireAdmin)
	adminGroup.POST("", h.CreateQuote)
	adminGroup.GET("", h.GetAdminQuotes)
	adminGroup.PUT("/:id", h.UpdateQuote)
	adminGroup.DELETE("/:id", h.DeleteQuote)

	// Public Routes
	publicGroup := e.Group("/api/v1/quotes", authMW)
	publicGroup.GET("", h.GetPublicQuotes)
}
