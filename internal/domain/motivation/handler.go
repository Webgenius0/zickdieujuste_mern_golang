package motivation

import (
	"errors"
	"net/http"

	"gotickets/internal/domain/motivation/dto"
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

// FindAll godoc
// @Summary      Get all motivations
// @Description  Retrieves a list of all motivation items.
// @Tags         Motivation
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.MotivationResponse
// @Router       /api/v1/motivations [get]
func (h *Handler) FindAll(c *echo.Context) error {
	res, err := h.svc.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to retrieve motivations", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// Create godoc
// @Summary      Create a motivation
// @Description  Admin endpoint to create a new motivation.
// @Tags         Motivation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateMotivationReq  true  "Motivation payload"
// @Success      201      {object}  dto.MotivationResponse
// @Failure      400      {object}  httpresponse.Error
// @Router       /api/v1/admin/motivations [post]
func (h *Handler) Create(c *echo.Context) error {
	var req dto.CreateMotivationReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	res, err := h.svc.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create motivation", err.Error()))
	}
	return c.JSON(http.StatusCreated, res)
}

// Update godoc
// @Summary      Update a motivation
// @Description  Admin endpoint to update an existing motivation.
// @Tags         Motivation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                   true  "Motivation ID"
// @Param        request  body      dto.UpdateMotivationReq  true  "Motivation payload"
// @Success      200      {object}  dto.MotivationResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Router       /api/v1/admin/motivations/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	var req dto.UpdateMotivationReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	res, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, ErrMotivationNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Motivation not found", ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update motivation", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary      Delete a motivation
// @Description  Admin endpoint to delete an existing motivation.
// @Tags         Motivation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Motivation ID"
// @Success      204  "No Content"
// @Failure      400  {object}  httpresponse.Error
// @Router       /api/v1/admin/motivations/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	if err := h.svc.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete motivation", err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}
