package server

import (
	"fmt"

	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/domain/content"
	"gotickets/internal/domain/schedule"
	"gotickets/internal/domain/subscription"
	"gotickets/internal/domain/user"
	"gotickets/internal/upload"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

// Start initializes and runs the HTTP server.
func Start(db *gorm.DB, cfg *config.Config, uploader upload.Uploader) {
	migrate(db)

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	// System
	e.GET("/", WelcomeHandler(cfg))
	e.GET("/health", HealthCheckHandler(db, cfg))
	RegisterSwagger(e)

	// Shared services
	jwtSvc := auth.NewJWTService(cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.JwtAccessExpiry, cfg.JwtRefreshExpiry)
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtSvc, uploader, nil)

	// Domain routes
	user.RegisterRoutes(e, db, cfg, uploader)
	content.RegisterRoutes(e, db, jwtSvc, userSvc)
	schedule.RegisterRoutes(e, db, jwtSvc, userSvc)
	subscription.RegisterRoutes(e, db, jwtSvc, userSvc, userRepo)

	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("\033[1;32m🚀 Server running on http://localhost:%s\033[0m\n", cfg.Port)
	if err := e.Start(addr); err != nil {
		e.Logger.Error("server stopped", "error", err)
	}
}

// migrate runs AutoMigrate for all domain entities.
func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		// user
		&user.User{},
		&user.RefreshToken{},
		&user.OTP{},
		&user.DeviceToken{},
		// content
		&content.Content{},
		&content.ContentAudience{},
		&content.RelatedContentJoin{},
		// schedule
		&schedule.UserSchedule{},
		// subscription
		&subscription.SubscriptionPlan{},
		&subscription.Subscription{},
	); err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}
}
