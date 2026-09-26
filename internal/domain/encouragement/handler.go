package encouragement

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gotickets/internal/httpresponse"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Create godoc
// @Summary      Create Encouragement
// @Description  Create a new encouragement (Admin)
// @Tags         admin, encouragements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateEncouragementReq  true  "Encouragement Data"
// @Success      201      {object}  EncouragementResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/encouragements [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateEncouragementReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation error", err.Error()))
	}

	res, err := h.svc.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create encouragement", err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

// GetAll godoc
// @Summary      Get All Encouragements
// @Description  Get a paginated list of encouragements
// @Tags         encouragements, admin
// @Produce      json
// @Param        page     query     int  false  "Page number"
// @Param        limit    query     int  false  "Items per page"
// @Success      200      {object}  PaginatedEncouragementResponse
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/encouragements [get]
// @Router       /api/v1/admin/encouragements [get]
func (h *Handler) GetAll(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	res, err := h.svc.GetAll(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch encouragements", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary      Get Encouragement by ID
// @Description  Get a specific encouragement by its ID
// @Tags         admin, encouragements
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Encouragement ID"
// @Success      200  {object}  EncouragementResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/encouragements/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	res, err := h.svc.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch encouragement", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// Update godoc
// @Summary      Update Encouragement
// @Description  Update an existing encouragement (Admin)
// @Tags         admin, encouragements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                   true  "Encouragement ID"
// @Param        request  body      UpdateEncouragementReq   true  "Updated Encouragement Data"
// @Success      200      {object}  EncouragementResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/encouragements/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	var req UpdateEncouragementReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation error", err.Error()))
	}

	res, err := h.svc.Update(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update encouragement", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary      Delete Encouragement
// @Description  Delete an encouragement by ID (Admin)
// @Tags         admin, encouragements
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Encouragement ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httpresponse.Error
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/encouragements/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	if err := h.svc.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete encouragement", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "encouragement deleted successfully"})
}
