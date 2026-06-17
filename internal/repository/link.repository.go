package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shortlink-backend/internal/model"
)

var ErrDuplicateSlug = errors.New("slug already used")
var ErrDuplicateOriginalURL = errors.New("original URL already used")

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (lr *LinkRepository) CreateLink(ctx context.Context, userId int, originalURL, slug string) (model.Link, error) {
	sql := `
		INSERT INTO links
		(user_id, original_url, slug)
		VALUES
		($1, $2, $3)
		RETURNING id, original_url, slug;
	`
	args := []any{userId, originalURL, slug}

	var link model.Link
	if err := lr.db.QueryRow(ctx, sql, args...).Scan(&link.ID, &link.OriginalURL, &link.Slug); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "links_slug_key" {
			return model.Link{}, ErrDuplicateSlug
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "links_original_url_key" {
			return model.Link{}, ErrDuplicateOriginalURL
		}
		return model.Link{}, err
	}

	return link, nil
}
