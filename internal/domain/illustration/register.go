package illustration

import (
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, authMW echo.MiddlewareFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	v1 := e.Group("/api/v1")

	// Public (Mobile) routes
	v1.GET("/illustrations", handler.GetAll)
	v1.GET("/illustrations/:id", handler.GetByID)

	// Admin routes
	admin := v1.Group("/admin/illustrations", authMW, middlewares.RequireAdmin)

	admin.POST("", handler.Create)
	admin.GET("", handler.GetAll)
	admin.GET("/:id", handler.GetByID)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}
