package content

import (
	"gotickets/internal/auth"
	"gotickets/internal/domain/user"
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes wires content domain routes onto the Echo router.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtSvc auth.JWTService, userSvc user.Service) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc, userSvc)

	authMW := middlewares.AuthMiddleware(jwtSvc)

	g := e.Group("/api/v1/content", authMW)
	g.GET("", h.ListContent)
	g.GET("/daily-quote", h.GetDailyQuote)
	g.GET("/:id", h.GetContentByID)
	g.GET("/:id/related", h.GetRelatedContent)
}
