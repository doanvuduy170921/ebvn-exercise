package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// urlRepo implements UrlRepo using Redis as the storage backend.
type urlRepo struct {
	r *redis.Client
}

// NewUrlRepo creates a new urlRepo with the given Redis client.
func NewUrlRepo(rdb *redis.Client) UrlRepo {
	return &urlRepo{r: rdb}
}

// Save stores a key-value pair in Redis with the given expiration time in seconds.
func (u *urlRepo) Save(ctx context.Context, key, value string, exp int) error {
	return u.r.Set(ctx, key, value, time.Duration(exp)*time.Second).Err()
}

// Get retrieves the value associated with the given key from Redis.
func (u *urlRepo) Get(ctx context.Context, key string) (string, error) {
	return u.r.Get(ctx, key).Result()
}
