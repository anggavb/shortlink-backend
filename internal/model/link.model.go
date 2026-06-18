package model

import "time"

type Link struct {
	ID          int64     `db:"id"`
	OriginalURL string    `db:"original_url"`
	Slug        string    `db:"slug"`
	CreatedAt   time.Time `db:"created_at"`
	ClickCount  int64     `db:"click_count"`
}
