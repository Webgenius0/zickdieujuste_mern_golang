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

// CreateCategory godoc
// @Summary      Create a category
// @Description  Creates a new prayer category. Requires ADMIN role.
// @Tags         Prayer Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateCategoryRequest  true  "Category data"
// @Success      201      {object}  dto.CategoryResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/categories [post]
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

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Returns a list of all prayer categories.
// @Tags         Prayer Categories
// @Produce      json
// @Success      200  {array}   dto.CategoryResponse
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/categories [get]
func (h *Handler) GetAllCategories(c *echo.Context) error {
	resp, err := h.svc.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError(http.StatusInternalServerError, "Failed to fetch categories", err.Error()))
	}
	return c.JSON(http.StatusOK, resp)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Returns a single category by its ID.
// @Tags         Prayer Categories
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  dto.CategoryResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/categories/{id} [get]
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

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Updates an existing prayer category. Requires ADMIN role.
// @Tags         Prayer Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                      true  "Category ID"
// @Param        request  body      dto.UpdateCategoryRequest   true  "Category data"
// @Success      200      {object}  dto.CategoryResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/categories/{id} [put]
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

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Deletes a prayer category and its subcategories. Requires ADMIN role.
// @Tags         Prayer Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Category ID"
// @Success      204  "No Content"
// @Failure      400  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/categories/{id} [delete]
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

// CreateSubCategory godoc
// @Summary      Create a subcategory
// @Description  Creates a new prayer subcategory. Requires ADMIN role.
// @Tags         Prayer SubCategories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateSubCategoryRequest  true  "SubCategory data"
// @Success      201      {object}  dto.SubCategoryResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/subcategories [post]
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

// GetSubCategoriesByCategoryID godoc
// @Summary      Get subcategories by category ID
// @Description  Returns a list of subcategories belonging to a specific category.
// @Tags         Prayer SubCategories
// @Produce      json
// @Param        categoryId  query     string  true  "Category ID"
// @Success      200         {array}   dto.SubCategoryResponse
// @Failure      400         {object}  httpresponse.Error
// @Failure      500         {object}  httpresponse.Error
// @Router       /api/v1/admin/subcategories [get]
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

// GetSubCategoryByID godoc
// @Summary      Get subcategory by ID
// @Description  Returns a single subcategory by its ID.
// @Tags         Prayer SubCategories
// @Produce      json
// @Param        id   path      string  true  "SubCategory ID"
// @Success      200  {object}  dto.SubCategoryResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/subcategories/{id} [get]
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

// UpdateSubCategory godoc
// @Summary      Update a subcategory
// @Description  Updates an existing prayer subcategory. Requires ADMIN role.
// @Tags         Prayer SubCategories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                         true  "SubCategory ID"
// @Param        request  body      dto.UpdateSubCategoryRequest   true  "SubCategory data"
// @Success      200      {object}  dto.SubCategoryResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/subcategories/{id} [put]
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

// DeleteSubCategory godoc
// @Summary      Delete a subcategory
// @Description  Deletes a prayer subcategory. Requires ADMIN role.
// @Tags         Prayer SubCategories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "SubCategory ID"
// @Success      204  "No Content"
// @Failure      400  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/subcategories/{id} [delete]
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

// CreatePrayer godoc
// @Summary      Create a prayer
// @Description  Creates a new prayer. Requires ADMIN role.
// @Tags         Prayers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreatePrayerRequest  true  "Prayer data"
// @Success      201      {object}  dto.PrayerResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/prayers [post]
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

// GetAllPrayers godoc
// @Summary      Get all prayers
// @Description  Returns a paginated list of prayers with optional filtering.
// @Tags         Prayers
// @Produce      json
// @Param        page            query     int     false  "Page number (default 1)"
// @Param        limit           query     int     false  "Items per page (default 10)"
// @Param        targetAudience  query     string  false  "Filter by target audience"
// @Param        categoryId      query     string  false  "Filter by category ID"
// @Param        ageGroup        query     string  false  "Filter by age group"
// @Param        mediaType       query     string  false  "Filter by media type"
// @Success      200             {object}  dto.PaginatedPrayerResponse
// @Failure      500             {object}  httpresponse.Error
// @Router       /api/v1/admin/prayers [get]
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

// GetPrayerByID godoc
// @Summary      Get prayer by ID
// @Description  Returns a single prayer by its ID.
// @Tags         Prayers
// @Produce      json
// @Param        id   path      string  true  "Prayer ID"
// @Success      200  {object}  dto.PrayerResponse
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/prayers/{id} [get]
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

// UpdatePrayer godoc
// @Summary      Update a prayer
// @Description  Updates an existing prayer. Requires ADMIN role.
// @Tags         Prayers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Prayer ID"
// @Param        request  body      dto.UpdatePrayerRequest   true  "Prayer data"
// @Success      200      {object}  dto.PrayerResponse
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/admin/prayers/{id} [put]
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

// DeletePrayer godoc
// @Summary      Delete a prayer
// @Description  Deletes a prayer by ID and removes its assets from Cloudinary. Requires ADMIN role.
// @Tags         Prayers
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Prayer ID"
// @Success      204  "No Content"
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/admin/prayers/{id} [delete]
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
