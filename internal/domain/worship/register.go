package worship

import (
	"gotickets/internal/auth"
	"gotickets/internal/middlewares"
	"gotickets/internal/upload"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes sets up the Echo routes for the Worship domain.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService, uploader upload.Uploader) {
	repo := NewRepository(db)
	svc := NewService(repo, uploader)
	h := NewHandler(svc)

	// Public routes
	public := e.Group("/api/v1")
	public.GET("/worships", h.GetAll)
	public.GET("/worships/:id", h.GetByID)

	// Admin routes
	admin := e.Group("/api/v1/admin/worships")
	admin.Use(middlewares.AuthMiddleware(jwtService))
	admin.Use(middlewares.RequireAdmin)

	admin.POST("", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}
