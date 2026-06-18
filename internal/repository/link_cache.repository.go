package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrLinkCacheMiss = errors.New("link cache miss")

type LinkCacheRepository struct {
	rdb    *redis.Client
	prefix string
}

func NewLinkCacheRepository(rdb *redis.Client) *LinkCacheRepository {
	return &LinkCacheRepository{
		rdb:    rdb,
		prefix: strings.TrimRight(os.Getenv("RDB_PREFIX"), ":"),
	}
}

func (r *LinkCacheRepository) SaveOriginalURL(ctx context.Context, slug, originalURL string, ttl time.Duration) error {
	return r.rdb.Set(ctx, r.linkKey(slug), originalURL, ttl).Err()
}

func (r *LinkCacheRepository) GetOriginalURL(ctx context.Context, slug string) (string, error) {
	originalURL, err := r.rdb.Get(ctx, r.linkKey(slug)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrLinkCacheMiss
		}
		return "", err
	}

	return originalURL, nil
}

func (r *LinkCacheRepository) DeleteOriginalURL(ctx context.Context, slug string) error {
	return r.rdb.Del(ctx, r.linkKey(slug)).Err()
}

func (r *LinkCacheRepository) linkKey(slug string) string {
	return fmt.Sprintf("%s:url:%s", r.prefix, slug)
}
