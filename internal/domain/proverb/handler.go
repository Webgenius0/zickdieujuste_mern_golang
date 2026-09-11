package proverb

import (
	"net/http"
	"strconv"

	"gotickets/internal/domain/proverb/dto"
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

// GetAll godoc
// @Summary      Get all proverbs
// @Description  Returns a paginated list of proverbs, ordered by publish_date DESC.
// @Tags         Proverb
// @Produce      json
// @Param        page         query     int     false  "Page number (default 1)"
// @Param        limit        query     int     false  "Items per page (default 10, max 100)"
// @Success      200  {object}  dto.PaginatedProverbResponse
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/proverbs [get]
func (h *Handler) GetAll(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	resp, err := h.svc.GetAll(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch proverbs", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// GetByID godoc
// @Summary      Get proverb by ID
// @Description  Returns a single proverb by its ID.
// @Tags         Proverb
// @Produce      json
// @Param        id   path      string  true  "Proverb ID"
// @Success      200  {object}  dto.ProverbResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Router       /api/v1/proverbs/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	resp, err := h.svc.GetByID(id)
	if err != nil {
		if err.Error() == "proverb not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch proverb", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// Create godoc
// @Summary      Create a proverb
// @Description  Creates a new proverb. Requires ADMIN role.
// @Tags         Proverb
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateProverbReq  true  "Proverb data"
// @Success      201      {object}  dto.ProverbResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/proverbs [post]
func (h *Handler) Create(c *echo.Context) error {
	var req dto.CreateProverbReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create proverb", err.Error()))
	}
	return c.JSON(http.StatusCreated, resp)
}

// Update godoc
// @Summary      Update a proverb
// @Description  Updates an existing proverb. Requires ADMIN role.
// @Tags         Proverb
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string            true  "Proverb ID"
// @Param        request  body      dto.UpdateProverbReq  true  "Proverb data"
// @Success      200      {object}  dto.ProverbResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/proverbs/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	var req dto.UpdateProverbReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.Update(id, req)
	if err != nil {
		if err.Error() == "proverb not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update proverb", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary      Delete a proverb
// @Description  Deletes a proverb by ID and removes its assets from Cloudinary. Requires ADMIN role.
// @Tags         Proverb
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Proverb ID"
// @Success      204  "No Content"
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/proverbs/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "proverb not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete proverb", err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}
