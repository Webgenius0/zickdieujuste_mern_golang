package home

import (
	"gotickets/internal/domain/prayer"
	"gotickets/internal/domain/quote"
	"gotickets/internal/domain/schedule"
	"gotickets/internal/domain/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, authMW echo.MiddlewareFunc) {
	// Instantiate repositories from existing domains
	userRepo := user.NewRepository(db)
	scheduleRepo := schedule.NewRepository(db)
	prayerRepo := prayer.NewRepository(db)
	quoteRepo := quote.NewRepository(db)

	svc := NewService(userRepo, scheduleRepo, prayerRepo, quoteRepo)
	handler := NewHandler(svc)

	v1 := e.Group("/api/v1")
	
	// Protected home route
	v1.GET("/home", handler.GetHomeData, authMW)
}
