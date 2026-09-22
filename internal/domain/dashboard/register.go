package dashboard

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, authMW echo.MiddlewareFunc) {
	h := NewHandler(db)

	group := e.Group("/api/v1/admin/dashboard", authMW)
	group.GET("/overview", h.GetOverview)
}
