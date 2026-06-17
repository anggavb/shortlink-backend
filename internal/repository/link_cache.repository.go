package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func (r *LinkCacheRepository) SaveOriginalURL(ctx context.Context, userId int, slug, originalURL string, ttl time.Duration) error {
	return r.rdb.Set(ctx, r.linkKey(userId, slug), originalURL, ttl).Err()
}

func (r *LinkCacheRepository) linkKey(userId int, slug string) string {
	return fmt.Sprintf("%s:%d:%s", r.prefix, userId, slug)
}
