package prayer

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gotickets/internal/domain/prayer/dto"
	"gotickets/internal/httpresponse"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Categories

func (h *Handler) CreateCategory(c *echo.Context) error {
	var req dto.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}
	resp, err := h.svc.CreateCategory(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create category", err.Error()))
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetAllCategories(c *echo.Context) error {
	resp, err := h.svc.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch categories", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCategoryByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	resp, err := h.svc.GetCategoryByID(id)
	if err != nil {
		if err.Error() == "category not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch category", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateCategory(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	var req dto.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}
	resp, err := h.svc.UpdateCategory(id, req)
	if err != nil {
		if err.Error() == "category not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update category", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteCategory(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	if err := h.svc.DeleteCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete category", err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}

// SubCategories

func (h *Handler) CreateSubCategory(c *echo.Context) error {
	var req dto.CreateSubCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}
	resp, err := h.svc.CreateSubCategory(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create subcategory", err.Error()))
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetSubCategoriesByCategoryID(c *echo.Context) error {
	catID, err := uuid.Parse(c.QueryParam("categoryId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid categoryId", err.Error()))
	}
	resp, err := h.svc.GetSubCategoriesByCategoryID(catID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch subcategories", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSubCategoryByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	resp, err := h.svc.GetSubCategoryByID(id)
	if err != nil {
		if err.Error() == "subcategory not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch subcategory", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateSubCategory(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	var req dto.UpdateSubCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}
	resp, err := h.svc.UpdateSubCategory(id, req)
	if err != nil {
		if err.Error() == "subcategory not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update subcategory", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteSubCategory(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	if err := h.svc.DeleteSubCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete subcategory", err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}

// Prayers

func (h *Handler) CreatePrayer(c *echo.Context) error {
	var req dto.CreatePrayerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}
	resp, err := h.svc.CreatePrayer(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to create prayer", err.Error()))
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetAllPrayers(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	filters := map[string]interface{}{}
	if t := c.QueryParam("targetAudience"); t != "" {
		filters["targetAudience"] = t
	}
	if cat := c.QueryParam("categoryId"); cat != "" {
		filters["categoryId"] = cat
	}
	if a := c.QueryParam("ageGroup"); a != "" {
		filters["ageGroup"] = a
	}
	if m := c.QueryParam("mediaType"); m != "" {
		filters["mediaType"] = m
	}

	resp, err := h.svc.GetAllPrayers(page, limit, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch prayers", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPrayerByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	resp, err := h.svc.GetPrayerByID(id)
	if err != nil {
		if err.Error() == "prayer not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch prayer", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdatePrayer(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	var req dto.UpdatePrayerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request body", err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}
	resp, err := h.svc.UpdatePrayer(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "prayer not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to update prayer", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeletePrayer(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid ID format", err.Error()))
	}
	if err := h.svc.DeletePrayer(c.Request().Context(), id); err != nil {
		if err.Error() == "prayer not found" {
			return c.JSON(http.StatusNotFound, httpresponse.NewError(http.StatusNotFound, err.Error(), ""))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to delete prayer", err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}
