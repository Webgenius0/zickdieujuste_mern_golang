package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/middlewares"
	"gotickets/internal/upload"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes wires all user-domain routes onto the Echo router.
// It instantiates all dependencies internally and applies the auth middleware where required.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, uploader upload.Uploader) {
	repo := NewRepository(db)
	jwtSvc := auth.NewJWTService(cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.JwtAccessExpiry, cfg.JwtRefreshExpiry)
	svc := NewService(repo, jwtSvc, uploader)
	h := NewHandler(svc, uploader)

	authMW := middlewares.AuthMiddleware(jwtSvc)

	// Auth routes — no middleware
	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/refresh", h.Refresh)
	authGroup.POST("/logout", h.Logout, authMW)
	authGroup.POST("/forgot-password", h.ForgotPassword)
	authGroup.POST("/reset-password", h.ResetPassword)

	// Profile routes — require auth
	userGroup := e.Group("/api/v1/users", authMW)
	userGroup.GET("/me", h.GetMe)
	userGroup.PUT("/me", h.UpdateMe)
	userGroup.PUT("/me/password", h.ChangePassword)
	userGroup.DELETE("/me", h.DeleteMe)
	userGroup.POST("/me/avatar", h.UploadAvatar)

	// Device routes — require auth
	deviceGroup := e.Group("/api/v1/devices", authMW)
	deviceGroup.POST("", h.RegisterDevice)
}
