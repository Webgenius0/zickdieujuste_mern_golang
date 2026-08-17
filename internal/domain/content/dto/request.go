package dto

// ListContentRequest holds query params for GET /api/v1/content.
type ListContentRequest struct {
	Type        string `query:"type"`
	SubType     string `query:"sub_type"`
	Audience    string `query:"audience"`
	CategoryTag string `query:"category_tag"`
	Q           string `query:"q"` // Full-text search query
	Page        int    `query:"page"`
	PageSize    int    `query:"page_size"`
}
