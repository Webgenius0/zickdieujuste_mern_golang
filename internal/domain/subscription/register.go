package subscription

import (
	"gotickets/internal/auth"
	"gotickets/internal/domain/user"
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes wires subscription domain routes onto the Echo router.
// It also seeds the subscription_plans table on startup if empty.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtSvc auth.JWTService, userSvc user.Service, userRepo user.Repository) {
	repo := NewRepository(db)
	verifier := NewNoopVerifier()
	svc := NewService(repo, verifier, userRepo)
	h := NewHandler(svc, userSvc)

	// Seed plans once on startup if the table is empty
	if err := SeedPlans(repo); err != nil {
		// Non-fatal: log and continue
		e.Logger.Warn("Failed to seed subscription plans: ", err)
	}

	authMW := middlewares.AuthMiddleware(jwtSvc)

	// Plans — requires auth
	subsGroup := e.Group("/api/v1/subscriptions", authMW)
	subsGroup.GET("/plans", h.ListPlans)
	subsGroup.POST("/verify", h.VerifyReceipt)

	// Webhook — no auth middleware (store-signed)
	e.POST("/api/v1/subscriptions/webhook", h.HandleWebhook)
}
