package cms

import (
	"net/http"

	"gotickets/internal/domain/cms/dto"
	"gotickets/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Admin APIs

func (h *Handler) GetAdminPages(c *echo.Context) error {
	res, err := h.svc.GetAllPages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch pages", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAdminPageBySlug(c *echo.Context) error {
	slug := c.Param("slug")
	res, err := h.svc.GetPageBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Page not found", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateAdminPage(c *echo.Context) error {
	slug := c.Param("slug")

	var req dto.UpdateCMSPageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	res, err := h.svc.UpdatePage(slug, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update page", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// Public APIs

func (h *Handler) GetPublicPage(c *echo.Context) error {
	slug := c.Param("slug")
	res, err := h.svc.GetPageBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Page not found", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}
