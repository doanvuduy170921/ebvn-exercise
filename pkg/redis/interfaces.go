package redis

import (
	"context"
	"time"
)

//go:generate mockery --name RedisClient --filename=mock.go --outpkg redis
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
}
