package home

import (
	"net/http"

	"gotickets/internal/auth"
	"gotickets/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GetHomeData godoc
// @Summary      Get Home Data
// @Description  Get aggregated data for the mobile home screen (User info, Schedule, Prayers, Quote)
// @Tags         home
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  HomeResponse
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/home [get]
func (h *Handler) GetHomeData(c *echo.Context) error {
	claims, ok := c.Get("user").(*auth.JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Unauthorized", "invalid token claims"))
	}

	res, err := h.svc.GetHomeData(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to load home data", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}
