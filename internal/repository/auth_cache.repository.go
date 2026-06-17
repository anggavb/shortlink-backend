package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shortlink-backend/internal/model"
)

type AuthCacheRepository struct {
	rdb    *redis.Client
	prefix string
}

func NewAuthCacheRepository(rdb *redis.Client) *AuthCacheRepository {
	return &AuthCacheRepository{
		rdb:    rdb,
		prefix: strings.TrimRight(os.Getenv("RDB_PREFIX"), ":"),
	}
}

func (r *AuthCacheRepository) SaveToken(ctx context.Context, user model.User, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("token already expired")
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, r.tokenKey(user.Id), userJSON, ttl).Err()
}

func (r *AuthCacheRepository) DeleteToken(ctx context.Context, userId int) error {
	return r.rdb.Del(ctx, r.tokenKey(userId)).Err()
}

func (r *AuthCacheRepository) tokenKey(userID int) string {
	return fmt.Sprintf("%s:user:%d", r.prefix, userID)
}
