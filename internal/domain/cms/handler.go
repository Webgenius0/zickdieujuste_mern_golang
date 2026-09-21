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

// CreateAdminPage godoc
// @Summary      Create CMS Page
// @Description  Create a new CMS page (e.g. Privacy Policy, Terms of Service)
// @Tags         admin, cms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateCMSPageRequest  true  "CMS Page Data"
// @Success      201      {object}  dto.CMSPageResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/cms/pages [post]
func (h *Handler) CreateAdminPage(c *echo.Context) error {
	var req dto.CreateCMSPageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	res, err := h.svc.CreatePage(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create page", err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

// GetAdminPages godoc
// @Summary      Get all CMS Pages
// @Description  Retrieve all CMS pages (Admin)
// @Tags         admin, cms
// @Produce      json
// @Security     BearerAuth
// @Success      200      {array}   dto.CMSPageResponse
// @Failure      401      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/cms/pages [get]
func (h *Handler) GetAdminPages(c *echo.Context) error {
	res, err := h.svc.GetAllPages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch pages", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// GetAdminPageBySlug godoc
// @Summary      Get CMS Page by Slug (Admin)
// @Description  Retrieve a specific CMS page by its slug
// @Tags         admin, cms
// @Produce      json
// @Security     BearerAuth
// @Param        slug     path      string  true  "CMS Page Slug"
// @Success      200      {object}  dto.CMSPageResponse
// @Failure      401      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Router       /api/v1/admin/cms/pages/{slug} [get]
func (h *Handler) GetAdminPageBySlug(c *echo.Context) error {
	slug := c.Param("slug")
	res, err := h.svc.GetPageBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Page not found", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateAdminPage godoc
// @Summary      Update CMS Page
// @Description  Update an existing CMS page and its sections
// @Tags         admin, cms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug     path      string                    true  "CMS Page Slug"
// @Param        request  body      dto.UpdateCMSPageRequest  true  "CMS Page Data"
// @Success      200      {object}  dto.CMSPageResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      401      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/cms/pages/{slug} [put]
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

// GetPrivacyPolicy godoc
// @Summary      Get Privacy Policy
// @Description  Retrieve the Privacy Policy page
// @Tags         cms
// @Produce      json
// @Success      200      {object}  dto.CMSPageResponse
// @Failure      404      {object}  httpresponse.Error
// @Router       /api/v1/cms/pages/privacy-policy [get]
func (h *Handler) GetPrivacyPolicy(c *echo.Context) error {
	res, err := h.svc.GetPageBySlug("privacy-policy")
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, "Privacy Policy not found", err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}
