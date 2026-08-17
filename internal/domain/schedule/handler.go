package schedule

import (
	"net/http"

	"gotickets/internal/auth"
	"gotickets/internal/domain/schedule/dto"
	"gotickets/internal/domain/user"
	"gotickets/internal/httpresponse"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// Handler holds the Echo HTTP handlers for the schedule domain.
type Handler struct {
	svc     Service
	userSvc user.Service
}

// NewHandler creates a new schedule Handler.
func NewHandler(svc Service, userSvc user.Service) *Handler {
	return &Handler{svc: svc, userSvc: userSvc}
}

// GetMySchedule godoc
// @Summary      Get current user's prayer schedule
// @Description  Returns the authenticated user's prayer schedule. Creates a default schedule on first access (05:00 morning, 21:00 night, UTC, push enabled) rather than returning 404.
// @Tags         Schedule
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ScheduleResponse
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/schedules/me [get]
func (h *Handler) GetMySchedule(c *echo.Context) error {
	userID, err := resolveUserID(c, h.userSvc)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Unauthorized", err.Error()))
	}

	resp, err := h.svc.GetOrCreateSchedule(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch schedule", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// UpdateMySchedule godoc
// @Summary      Update prayer schedule
// @Description  Updates morning_prayer_time, night_prayer_time, timezone (must be valid IANA string), and push_enabled for the authenticated user.
// @Tags         Schedule
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.UpdateScheduleRequest  true  "Schedule update payload"
// @Success      200      {object}  dto.ScheduleResponse
// @Failure      400      {object}  httpresponse.Error  "Validation error or invalid timezone"
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/schedules/me [put]
func (h *Handler) UpdateMySchedule(c *echo.Context) error {
	userID, err := resolveUserID(c, h.userSvc)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Unauthorized", err.Error()))
	}

	var req dto.UpdateScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.UpdateSchedule(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, err.Error(), ""))
	}
	return c.JSON(http.StatusOK, resp)
}

// resolveUserID extracts the email from JWT claims and resolves the UUID via user service.
func resolveUserID(c *echo.Context, userSvc user.Service) (uuid.UUID, error) {
	claims, ok := c.Get("user").(*auth.JwtCustomClaims)
	if !ok || claims == nil {
		return uuid.Nil, errUnauthorized()
	}
	id, err := userSvc.GetUserIDByEmail(claims.Email)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func errUnauthorized() error {
	return echo.ErrUnauthorized
}
