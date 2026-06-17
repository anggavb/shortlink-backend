package model

type Link struct {
	ID          int64  `db:"id"`
	OriginalURL string `db:"original_url"`
	Slug        string `db:"slug"`
}
