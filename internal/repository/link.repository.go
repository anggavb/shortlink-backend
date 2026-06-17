package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shortlink-backend/internal/model"
)

var ErrDuplicateSlug = errors.New("slug already used")
var ErrDuplicateOriginalURL = errors.New("original URL already used")
var ErrLinkNotFound = errors.New("link not found")

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

func (lr *LinkRepository) CountLinksByUser(ctx context.Context, userId int) (int64, error) {
	sql := `
		SELECT COUNT(id)
		FROM links
		WHERE user_id = $1
			AND deleted_at IS NULL;
	`

	var total int64
	if err := lr.db.QueryRow(ctx, sql, userId).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (lr *LinkRepository) ListLinksByUser(ctx context.Context, userId, limit, offset int) ([]model.Link, error) {
	sql := `
		SELECT id, original_url, slug
		FROM links
		WHERE user_id = $1
			AND deleted_at IS NULL
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3;
	`
	args := []any{userId, limit, offset}

	rows, err := lr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]model.Link, 0)
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(&link.ID, &link.OriginalURL, &link.Slug); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (lr *LinkRepository) SoftDeleteLink(ctx context.Context, userId int, linkId int64) (string, error) {
	sql := `
		UPDATE links
		SET deleted_at = now(),
			updated_at = now()
		WHERE id = $1
			AND user_id = $2
			AND deleted_at IS NULL
		RETURNING slug;
	`
	args := []any{linkId, userId}

	var slug string
	if err := lr.db.QueryRow(ctx, sql, args...).Scan(&slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}

	return slug, nil
}
