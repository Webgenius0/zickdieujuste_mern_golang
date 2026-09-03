package server

import (
	"errors"
	"fmt"
	"strings"

	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/domain/content"
	"gotickets/internal/domain/media"
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
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			var errMsgs []string
			for _, e := range valErrs {
				errMsgs = append(errMsgs, humanizeValidationError(e))
			}
			return fmt.Errorf("%s", strings.Join(errMsgs, "; "))
		}
		return err
	}
	return nil
}

func humanizeValidationError(e validator.FieldError) string {
	field := e.Field()
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		switch e.Kind().String() {
		case "string":
			return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
		default:
			return fmt.Sprintf("%s must be at least %s", field, e.Param())
		}
	case "max":
		switch e.Kind().String() {
		case "string":
			return fmt.Sprintf("%s must be no more than %s characters long", field, e.Param())
		default:
			return fmt.Sprintf("%s must be no more than %s", field, e.Param())
		}
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, strings.ReplaceAll(e.Param(), " ", ", "))
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func Start(db *gorm.DB, cfg *config.Config, uploader upload.Uploader) {
	migrate(db)
	user.SeedAdmin(db, cfg.AdminEmail, cfg.AdminPassword)

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/", WelcomeHandler(cfg))
	e.GET("/health", HealthCheckHandler(db, cfg))
	RegisterSwagger(e)

	jwtSvc := auth.NewJWTService(cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.JwtAccessExpiry, cfg.JwtRefreshExpiry)
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtSvc, uploader, nil, nil)

	user.RegisterRoutes(e, db, cfg, uploader)
	content.RegisterRoutes(e, db, jwtSvc, userSvc)
	schedule.RegisterRoutes(e, db, jwtSvc, userSvc)
	subscription.RegisterRoutes(e, db, jwtSvc, userSvc, userRepo)
	media.RegisterRoutes(e, jwtSvc, uploader)

	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("\033[1;32m🚀 Server running on http://localhost:%s\033[0m\n", cfg.Port)
	if err := e.Start(addr); err != nil {
		e.Logger.Error("server stopped", "error", err)
	}
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&user.User{},
		&user.RefreshToken{},
		&user.OTP{},
		&user.DeviceToken{},
		&content.Content{},
		&content.ContentAudience{},
		&content.RelatedContentJoin{},
		&schedule.UserSchedule{},
		&subscription.SubscriptionPlan{},
		&subscription.Subscription{},
	); err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}
}
