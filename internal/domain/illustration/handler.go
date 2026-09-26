package illustration

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
// @Summary      Create Illustration
// @Description  Create a new illustration (Admin)
// @Tags         admin, illustrations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateIllustrationReq  true  "Illustration Data"
// @Success      201      {object}  IllustrationResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/illustrations [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateIllustrationReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation error", err.Error()))
	}

	res, err := h.svc.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create illustration", err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

// GetAll godoc
// @Summary      Get All Illustrations
// @Description  Get a paginated list of illustrations
// @Tags         illustrations, admin
// @Produce      json
// @Param        page     query     int  false  "Page number"
// @Param        limit    query     int  false  "Items per page"
// @Success      200      {object}  PaginatedIllustrationResponse
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/illustrations [get]
// @Router       /api/v1/admin/illustrations [get]
func (h *Handler) GetAll(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	res, err := h.svc.GetAll(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch illustrations", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary      Get Illustration by ID
// @Description  Get a specific illustration by its ID
// @Tags         admin, illustrations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Illustration ID"
// @Success      200  {object}  IllustrationResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/illustrations/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	res, err := h.svc.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch illustration", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// Update godoc
// @Summary      Update Illustration
// @Description  Update an existing illustration (Admin)
// @Tags         admin, illustrations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                 true  "Illustration ID"
// @Param        request  body      UpdateIllustrationReq  true  "Updated Illustration Data"
// @Success      200      {object}  IllustrationResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/illustrations/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	var req UpdateIllustrationReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation error", err.Error()))
	}

	res, err := h.svc.Update(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update illustration", err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary      Delete Illustration
// @Description  Delete an illustration by ID (Admin)
// @Tags         admin, illustrations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Illustration ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httpresponse.Error
// @Failure      401  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/illustrations/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	if err := h.svc.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete illustration", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "illustration deleted successfully"})
}
