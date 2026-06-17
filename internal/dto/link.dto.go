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

type GetLinkResponse struct {
	OriginalURL string `json:"original_url"`
}
