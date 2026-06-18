package dto

import "time"

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

type LinkListItemResponse struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	Slug        string    `json:"slug"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int64     `json:"click_count"`
}

type ListLinksQuery struct {
	Page   *int    `form:"page" binding:"omitempty,gte=1"`
	Limit  *int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Search *string `form:"search" binding:"omitempty,max=255"`
}

type ListLinksResponse struct {
	Data []LinkListItemResponse `json:"data"`
	Meta PaginationMeta         `json:"meta"`
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

type LinkURI struct {
	ID string `uri:"id" json:"id" binding:"required,numeric"`
}
