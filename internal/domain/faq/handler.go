package faq

import (
	"gotickets/internal/domain/faq/dto"
	"gotickets/internal/httpresponse"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateFAQ(c *echo.Context) error {
	var req dto.CreateFAQRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	res, err := h.svc.CreateFAQ(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create FAQ", err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetAdminFAQs(c *echo.Context) error {
	res, err := h.svc.GetAdminFAQs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch FAQs", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetPublicFAQs(c *echo.Context) error {
	search := c.QueryParam("search")

	res, err := h.svc.GetPublicFAQs(search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch FAQs", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateFAQ(c *echo.Context) error {
	id := c.Param("id")

	var req dto.UpdateFAQRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	res, err := h.svc.UpdateFAQ(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update FAQ", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteFAQ(c *echo.Context) error {
	id := c.Param("id")

	if err := h.svc.DeleteFAQ(id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete FAQ", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "FAQ deleted successfully"})
}
