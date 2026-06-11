package repository

import (
	"context"
	redisPkg "lesson01-ebvn/pkg/redis"
	"time"
)

// urlRepo implements UrlRepo using Redis as the storage backend.
type urlRepo struct {
	r redisPkg.RedisClient
}

// NewUrlRepo creates a new urlRepo with the given Redis client.
func NewUrlRepo(rdb redisPkg.RedisClient) UrlRepo {
	return &urlRepo{r: rdb}
}

// Save stores a key-value pair in Redis with the given expiration time in seconds.
func (u *urlRepo) Save(ctx context.Context, key, value string, exp int) error {
	return u.r.Set(ctx, key, value, time.Duration(exp)*time.Second)
}

// Get retrieves the value associated with the given key from Redis.
func (u *urlRepo) Get(ctx context.Context, key string) (string, error) {
	return u.r.Get(ctx, key)
}
