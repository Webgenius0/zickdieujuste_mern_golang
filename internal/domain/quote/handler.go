package quote

import (
	"net/http"

	"gotickets/internal/httpresponse"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// CreateQuote godoc
// @Summary      Create Quote
// @Description  Create a new daily quote (Admin)
// @Tags         admin, quotes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateQuoteRequest  true  "Quote Data"
// @Success      201      {object}  QuoteResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/quotes [post]
func (h *Handler) CreateQuote(c *echo.Context) error {
	var req CreateQuoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	res, err := h.svc.CreateQuote(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create quote", err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

// GetAdminQuotes godoc
// @Summary      Get Admin Quotes
// @Description  Retrieve all quotes including future scheduled ones (Admin)
// @Tags         admin, quotes
// @Produce      json
// @Security     BearerAuth
// @Success      200      {array}   QuoteResponse
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/quotes [get]
func (h *Handler) GetAdminQuotes(c *echo.Context) error {
	res, err := h.svc.GetAllQuotes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch quotes", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateQuote godoc
// @Summary      Update Quote
// @Description  Update an existing daily quote (Admin)
// @Tags         admin, quotes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string              true  "Quote ID"
// @Param        request  body      UpdateQuoteRequest  true  "Quote Data"
// @Success      200      {object}  QuoteResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/quotes/{id} [put]
func (h *Handler) UpdateQuote(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid UUID format", err.Error()))
	}

	var req UpdateQuoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	res, err := h.svc.UpdateQuote(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update quote", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteQuote godoc
// @Summary      Delete Quote
// @Description  Delete a daily quote (Admin)
// @Tags         admin, quotes
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Quote ID"
// @Success      200      {object}  httpresponse.Error
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/quotes/{id} [delete]
func (h *Handler) DeleteQuote(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid UUID format", err.Error()))
	}

	if err := h.svc.DeleteQuote(id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete quote", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Quote deleted successfully"})
}

// GetPublicQuotes godoc
// @Summary      Get Public Quotes
// @Description  Retrieve published daily quotes for users
// @Tags         quotes
// @Produce      json
// @Success      200      {array}   QuoteResponse
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/quotes [get]
func (h *Handler) GetPublicQuotes(c *echo.Context) error {
	res, err := h.svc.GetPublishedQuotes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch quotes", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// GetPublicQuoteDetails godoc
// @Summary      Get Public Quote Details
// @Description  Retrieve a specific published daily quote by ID
// @Tags         quotes
// @Produce      json
// @Param        id       path      string  true  "Quote ID"
// @Success      200      {object}  QuoteResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/quotes/{id} [get]
func (h *Handler) GetPublicQuoteDetails(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid UUID format", err.Error()))
	}

	res, err := h.svc.GetPublishedQuoteByID(id)
	if err != nil {
		if err == echo.ErrNotFound {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Quote not found or not published", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch quote details", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}
