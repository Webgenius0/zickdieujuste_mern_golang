package server

import (
	"fmt"
	"net/http"
	"time"

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
// It runs AutoMigrate, seeds reference data, registers all domain routes, and starts listening.
func Start(db *gorm.DB, cfg *config.Config, uploader upload.Uploader) {
	// AutoMigrate all domain entities
	if err := db.AutoMigrate(
		// user domain
		&user.User{},
		&user.RefreshToken{},
		&user.OTP{},
		&user.DeviceToken{},
		// content domain
		&content.Content{},
		&content.ContentAudience{},
		&content.RelatedContentJoin{},
		// schedule domain
		&schedule.UserSchedule{},
		// subscription domain
		&subscription.SubscriptionPlan{},
		&subscription.Subscription{},
	); err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	// System routes
	e.GET("/", WelcomeHandler)
	e.GET("/health", HealthCheckHandler(db))

	// Documentation (Swagger UI at /swagger/index.html)
	RegisterSwagger(e)

	// Build shared JWT service — used by all domain route registrations
	jwtSvc := auth.NewJWTService(cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.JwtAccessExpiry, cfg.JwtRefreshExpiry)

	// Domain 1: User (auth + profile + devices)
	user.RegisterRoutes(e, db, cfg, uploader)

	// Build a user service reference for domains that need premium checks or user ID resolution
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtSvc, uploader)

	// Domain 2: Content
	content.RegisterRoutes(e, db, jwtSvc, userSvc)

	// Domain 3: Schedule
	schedule.RegisterRoutes(e, db, jwtSvc, userSvc)

	// Domain 4: Subscription
	subscription.RegisterRoutes(e, db, jwtSvc, userSvc, userRepo)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("\033[1;32m🚀 Server is running on http://localhost:%s\033[0m\n", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// System handlers
// ──────────────────────────────────────────────────────────────────────────────

// HealthCheckHandler godoc
// @Summary      Health Check
// @Description  Check the health status of the API and the database connection.
// @Tags         System
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthCheckHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		dbStatus := "up"
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}
		return c.JSON(http.StatusOK, HealthResponse{
			Status:    "up",
			Database:  dbStatus,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}
