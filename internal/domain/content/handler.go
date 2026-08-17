package content

import (
	"net/http"

	"gotickets/internal/auth"
	"gotickets/internal/domain/content/dto"
	"gotickets/internal/domain/user"
	"gotickets/internal/httpresponse"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// Handler holds the Echo HTTP handlers for the content domain.
type Handler struct {
	svc     Service
	userSvc user.Service
}

// NewHandler creates a new content Handler.
func NewHandler(svc Service, userSvc user.Service) *Handler {
	return &Handler{svc: svc, userSvc: userSvc}
}

// ListContent godoc
// @Summary      List content
// @Description  Returns a paginated, filterable list of content items. Supports `type`, `sub_type`, `audience`, `category_tag`, and `q` (full-text search) query parameters.
// @Tags         Content
// @Produce      json
// @Security     BearerAuth
// @Param        type          query     string  false  "Content type (PRAYER, MOTIVATION, WORSHIP, PROVERB, DAILY_QUOTE, ILLUSTRATION, ENCOURAGEMENT)"
// @Param        sub_type      query     string  false  "Sub-type (e.g. Morning, Night)"
// @Param        audience      query     string  false  "Audience (ALL, KIDS, TEENS)"
// @Param        category_tag  query     string  false  "Category tag (e.g. Thanksgiving, Intercession)"
// @Param        q             query     string  false  "Full-text search query"
// @Param        page          query     int     false  "Page number (default: 1)"
// @Param        page_size     query     int     false  "Items per page (default: 20, max: 100)"
// @Success      200           {object}  dto.ContentListResponse
// @Failure      401           {object}  httpresponse.Error
// @Failure      500           {object}  httpresponse.Error
// @Router       /api/v1/content [get]
func (h *Handler) ListContent(c *echo.Context) error {
	var filter dto.ListContentRequest
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid query parameters", err.Error()))
	}

	result, err := h.svc.List(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch content", err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}

// GetDailyQuote godoc
// @Summary      Get today's daily quote
// @Description  Returns the DAILY_QUOTE content item for the current date, filtered by published_at.
// @Tags         Content
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ContentDetail
// @Failure      401  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error  "No quote for today"
// @Router       /api/v1/content/daily-quote [get]
func (h *Handler) GetDailyQuote(c *echo.Context) error {
	isPremium := resolveIsPremium(c, h.userSvc)

	result, err := h.svc.GetDailyQuote(isPremium)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "No daily quote available for today", ""))
	}
	return c.JSON(http.StatusOK, result)
}

// GetContentByID godoc
// @Summary      Get content detail
// @Description  Returns full detail for a content item. If the content is premium and the user is not premium, returns HTTP 403 with an upgrade prompt.
// @Tags         Content
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Content UUID"
// @Success      200  {object}  dto.ContentDetail
// @Failure      401  {object}  httpresponse.Error
// @Failure      403  {object}  dto.PremiumGateResponse  "Premium content — upgrade required"
// @Failure      404  {object}  httpresponse.Error
// @Router       /api/v1/content/{id} [get]
func (h *Handler) GetContentByID(c *echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid content ID", err.Error()))
	}

	isPremium := resolveIsPremium(c, h.userSvc)

	result, err := h.svc.GetByID(id, isPremium)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Content not found", err.Error()))
	}

	// Premium gate: content is premium and user is not premium
	if result.IsPremium && !isPremium {
		return c.JSON(http.StatusForbidden, dto.PremiumGateResponse{
			Code:    http.StatusForbidden,
			Message: "This content requires a premium subscription.",
			Upgrade: "/api/v1/subscriptions/plans",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetRelatedContent godoc
// @Summary      Get related content
// @Description  Returns a list of content items related to the given content ID via the related_content join table.
// @Tags         Content
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Content UUID"
// @Success      200  {array}   dto.ContentSummary
// @Failure      401  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Router       /api/v1/content/{id}/related [get]
func (h *Handler) GetRelatedContent(c *echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid content ID", err.Error()))
	}

	result, err := h.svc.GetRelated(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Content not found", err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

func parseUUID(c *echo.Context, param string) (uuid.UUID, error) {
	return uuid.Parse(c.Param(param))
}

// resolveIsPremium checks JWT claims to get the user email, then looks up is_premium.
// Returns false (non-premium) on any error to fail safe.
func resolveIsPremium(c *echo.Context, userSvc user.Service) bool {
	claims, ok := c.Get("user").(*auth.JwtCustomClaims)
	if !ok || claims == nil {
		return false
	}
	profile, err := userSvc.GetProfileByEmail(claims.Email)
	if err != nil {
		return false
	}
	return profile.IsPremium
}
