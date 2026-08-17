package subscription

import (
	"net/http"

	"gotickets/internal/auth"
	"gotickets/internal/domain/subscription/dto"
	"gotickets/internal/domain/user"
	"gotickets/internal/httpresponse"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// Handler holds the Echo HTTP handlers for the subscription domain.
type Handler struct {
	svc     Service
	userSvc user.Service
}

// NewHandler creates a new subscription Handler.
func NewHandler(svc Service, userSvc user.Service) *Handler {
	return &Handler{svc: svc, userSvc: userSvc}
}

// ListPlans godoc
// @Summary      List subscription plans
// @Description  Returns all active subscription plans (Biannual, Annual, Friends & Family).
// @Tags         Subscriptions
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.PlanResponse
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/subscriptions/plans [get]
func (h *Handler) ListPlans(c *echo.Context) error {
	plans, err := h.svc.ListPlans()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch plans", err.Error()))
	}
	return c.JSON(http.StatusOK, plans)
}

// VerifyReceipt godoc
// @Summary      Verify purchase receipt
// @Description  Verifies an Apple or Google Play receipt, upserts a subscription record, and marks the user as premium. Receipt validation is currently stubbed — real integration pending.
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.VerifyReceiptRequest  true  "Store and receipt payload"
// @Success      200      {object}  dto.SubscriptionResponse
// @Failure      400      {object}  httpresponse.Error  "Validation error"
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/subscriptions/verify [post]
func (h *Handler) VerifyReceipt(c *echo.Context) error {
	userID, err := resolveUserID(c, h.userSvc)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Unauthorized", err.Error()))
	}

	var req dto.VerifyReceiptRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.VerifyReceipt(userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Receipt verification failed", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleWebhook godoc
// @Summary      Handle store webhook
// @Description  Receives renewal, cancellation, and refund events from Apple and Google. No auth middleware — store signature verification is applied instead (currently stubbed with a TODO).
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        request  body      dto.WebhookRequest  true  "Webhook payload from Apple or Google"
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/subscriptions/webhook [post]
func (h *Handler) HandleWebhook(c *echo.Context) error {
	// TODO: verify store signature before processing
	var req dto.WebhookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid webhook payload", err.Error()))
	}

	if err := h.svc.HandleWebhook(req); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Webhook processing failed", err.Error()))
	}
	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "Webhook received"})
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

func resolveUserID(c *echo.Context, userSvc user.Service) (uuid.UUID, error) {
	claims, ok := c.Get("user").(*auth.JwtCustomClaims)
	if !ok || claims == nil {
		return uuid.Nil, echo.ErrUnauthorized
	}
	return userSvc.GetUserIDByEmail(claims.Email)
}
