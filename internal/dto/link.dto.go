package dto

type CreateLinkRequest struct {
	OriginalURL string  `json:"original_url" binding:"required,url" example:"https://www.example.com/very/long/url"`
	Slug        *string `json:"slug,omitempty" binding:"omitempty,shortlink_slug" example:"abc123"`
}

type CreateLinkResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	Slug        string `json:"slug"`
	ShortURL    string `json:"short_url"`
}

type ListLinksQuery struct {
	Page  *int `form:"page" binding:"omitempty,gte=1"`
	Limit *int `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

type ListLinksResponse struct {
	Data []CreateLinkResponse `json:"data"`
	Meta PaginationMeta       `json:"meta"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type GetLinkResponse struct {
	OriginalURL string `json:"original_url"`
}
