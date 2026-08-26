package media

import (
	"gotickets/internal/auth"
	"gotickets/internal/middlewares"
	"gotickets/internal/upload"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts the media upload/delete routes onto the Echo router.
func RegisterRoutes(e *echo.Echo, jwtSvc auth.JWTService, uploader upload.Uploader) {
	h := NewHandler(uploader)
	authMW := middlewares.AuthMiddleware(jwtSvc)
	e.POST("/api/v1/upload", h.Upload, authMW)
	e.DELETE("/api/v1/upload", h.Delete, authMW)
}
