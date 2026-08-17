package schedule

import (
	"gotickets/internal/auth"
	"gotickets/internal/domain/user"
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes wires schedule domain routes onto the Echo router.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtSvc auth.JWTService, userSvc user.Service) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc, userSvc)

	authMW := middlewares.AuthMiddleware(jwtSvc)

	g := e.Group("/api/v1/schedules", authMW)
	g.GET("/me", h.GetMySchedule)
	g.PUT("/me", h.UpdateMySchedule)
}
