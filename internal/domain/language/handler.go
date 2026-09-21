package language

import (
	"errors"
	"net/http"
	"strconv"

	"gotickets/internal/domain/language/dto"
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

// CreateLanguage godoc
// @Summary      Create a new language
// @Description  Creates a new language for the application. Admin only.
// @Tags         Admin Languages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateLanguageReq  true  "Language data"
// @Success      201      {object}  dto.LanguageResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      403      {object}  httpresponse.Error
// @Failure      409      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/languages [post]
func (h *Handler) CreateLanguage(c *echo.Context) error {
	var req dto.CreateLanguageReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.CreateLanguage(req)
	if err != nil {
		if errors.Is(err, ErrLanguageConflict) {
			return c.JSON(http.StatusConflict, httpresponse.NewError(http.StatusConflict, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create language", err.Error()))
	}
	return c.JSON(http.StatusCreated, resp)
}

// GetAdminLanguages godoc
// @Summary      Get all languages (Admin)
// @Description  Retrieves a paginated list of all languages including inactive ones. Admin only.
// @Tags         Admin Languages
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number"  default(1)
// @Param        limit   query     int     false  "Items per page"  default(10)
// @Param        search  query     string  false  "Search by name or code"
// @Success      200     {object}  dto.PaginatedLanguageResponse
// @Failure      401     {object}  httpresponse.Error
// @Failure      403     {object}  httpresponse.Error
// @Failure      500     {object}  httpresponse.Error
// @Router       /api/v1/admin/languages [get]
func (h *Handler) GetAdminLanguages(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	search := c.QueryParam("search")

	resp, err := h.svc.GetLanguages(page, limit, search, false) // false = include inactive
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch languages", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// UpdateLanguage godoc
// @Summary      Update a language
// @Description  Updates a language by ID. Admin only.
// @Tags         Admin Languages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                 true  "Language ID"
// @Param        request  body      dto.UpdateLanguageReq  true  "Language data to update"
// @Success      200      {object}  dto.LanguageResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      403      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      409      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/languages/{id} [put]
func (h *Handler) UpdateLanguage(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	var req dto.UpdateLanguageReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	resp, err := h.svc.UpdateLanguage(id, req)
	if err != nil {
		if errors.Is(err, ErrLanguageNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		if errors.Is(err, ErrLanguageConflict) {
			return c.JSON(http.StatusConflict, httpresponse.NewError(http.StatusConflict, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update language", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteLanguage godoc
// @Summary      Delete a language
// @Description  Permanently deletes a language by ID. Admin only.
// @Tags         Admin Languages
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Language ID"
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      403      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/languages/{id} [delete]
func (h *Handler) DeleteLanguage(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}

	if err := h.svc.DeleteLanguage(id); err != nil {
		if errors.Is(err, ErrLanguageNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete language", err.Error()))
	}
	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "Language deleted successfully"})
}

// GetPublicLanguages godoc
// @Summary      Get active languages
// @Description  Retrieves a paginated list of all ACTIVE languages for the frontend/mobile app.
// @Tags         Languages
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number"  default(1)
// @Param        limit   query     int     false  "Items per page"  default(50)
// @Param        search  query     string  false  "Search by name or code"
// @Success      200     {object}  dto.PaginatedLanguageResponse
// @Failure      401     {object}  httpresponse.Error
// @Failure      500     {object}  httpresponse.Error
// @Router       /api/v1/languages [get]
func (h *Handler) GetPublicLanguages(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 50
	}
	search := c.QueryParam("search")

	resp, err := h.svc.GetLanguages(page, limit, search, true) // true = active only
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch languages", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}
