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

// CreateFAQ godoc
// @Summary      Create FAQ
// @Description  Create a new FAQ entry
// @Tags         admin, faqs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateFAQRequest  true  "FAQ Data"
// @Success      201      {object}  dto.FAQResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/faqs [post]
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

// GetAdminFAQs godoc
// @Summary      Get Admin FAQs
// @Description  Retrieve all FAQs (including inactive ones)
// @Tags         admin, faqs
// @Produce      json
// @Security     BearerAuth
// @Success      200      {array}   dto.FAQResponse
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/faqs [get]
func (h *Handler) GetAdminFAQs(c *echo.Context) error {
	res, err := h.svc.GetAdminFAQs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch FAQs", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// GetPublicFAQs godoc
// @Summary      Get Public FAQs
// @Description  Retrieve active FAQs, optionally filtered by a search query
// @Tags         faqs
// @Produce      json
// @Param        search  query      string  false  "Search query for FAQs"
// @Success      200     {array}    dto.FAQResponse
// @Failure      500     {object}   httpresponse.Error
// @Router       /api/v1/faqs [get]
func (h *Handler) GetPublicFAQs(c *echo.Context) error {
	search := c.QueryParam("search")

	res, err := h.svc.GetPublicFAQs(search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch FAQs", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateFAQ godoc
// @Summary      Update FAQ
// @Description  Update an existing FAQ entry
// @Tags         admin, faqs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "FAQ ID"
// @Param        request  body      dto.UpdateFAQRequest  true  "FAQ Data"
// @Success      200      {object}  dto.FAQResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/faqs/{id} [put]
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

// DeleteFAQ godoc
// @Summary      Delete FAQ
// @Description  Delete an FAQ entry by ID
// @Tags         admin, faqs
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "FAQ ID"
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/faqs/{id} [delete]
func (h *Handler) DeleteFAQ(c *echo.Context) error {
	id := c.Param("id")

	if err := h.svc.DeleteFAQ(id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete FAQ", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "FAQ deleted successfully"})
}
