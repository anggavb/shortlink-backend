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

func (r *AuthCacheRepository) SaveToken(ctx context.Context, tokenHash string, user model.User, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("token already expired")
	}

	userWithToken := model.UserWithToken{
		User:  user,
		Token: tokenHash,
	}

	userJSON, err := json.Marshal(userWithToken)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, r.tokenKey(user.Id, tokenHash), userJSON, ttl).Err()
}

func (r *AuthCacheRepository) tokenKey(userID int, tokenHash string) string {
	key := fmt.Sprintf("auth:token~%d:%s", userID, tokenHash)
	if r.prefix == "" {
		return key
	}

	return r.prefix + ":" + key
}
