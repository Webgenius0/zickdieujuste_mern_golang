package worship

import (
	"net/http"
	"strconv"

	"gotickets/internal/domain/worship/dto"
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
// @Summary      Get all worship tracks
// @Description  Returns a paginated list of worship tracks. Optionally filter by time_of_day.
// @Tags         Worship
// @Produce      json
// @Param        time_of_day  query     string  false  "Filter by Day or Night"
// @Param        page         query     int     false  "Page number (default 1)"
// @Param        limit        query     int     false  "Items per page (default 10, max 100)"
// @Success      200  {object}  dto.PaginatedWorshipResponse
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/worships [get]
func (h *Handler) GetAll(c *echo.Context) error {
	timeOfDay := c.QueryParam("time_of_day")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	resp, err := h.svc.GetAll(timeOfDay, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch worship tracks", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// GetByID godoc
// @Summary      Get worship track by ID
// @Description  Returns a single worship track by its ID.
// @Tags         Worship
// @Produce      json
// @Param        id   path      string  true  "Worship ID"
// @Success      200  {object}  dto.WorshipResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Router       /api/v1/worships/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	resp, err := h.svc.GetByID(id)
	if err != nil {
		if err.Error() == "worship track not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch worship track", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// Create godoc
// @Summary      Create a worship track
// @Description  Creates a new worship track. Requires ADMIN role.
// @Tags         Worship
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateWorshipReq  true  "Worship data"
// @Success      201      {object}  dto.WorshipResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/worships [post]
func (h *Handler) Create(c *echo.Context) error {
	var req dto.CreateWorshipReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create worship track", err.Error()))
	}
	return c.JSON(http.StatusCreated, resp)
}

// Update godoc
// @Summary      Update a worship track
// @Description  Updates an existing worship track. Requires ADMIN role.
// @Tags         Worship
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string            true  "Worship ID"
// @Param        request  body      dto.UpdateWorshipReq  true  "Worship data"
// @Success      200      {object}  dto.WorshipResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/worships/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	var req dto.UpdateWorshipReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.Update(id, req)
	if err != nil {
		if err.Error() == "worship track not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update worship track", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary      Delete a worship track
// @Description  Deletes a worship track by ID and removes its assets from Cloudinary. Requires ADMIN role.
// @Tags         Worship
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Worship ID"
// @Success      204  "No Content"
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/worships/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "worship track not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete worship track", err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}
