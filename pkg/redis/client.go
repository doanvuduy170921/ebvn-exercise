package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client using environment-based configuration.
// It pings Redis to verify the connection before returning.
func NewRedisClient() (RedisClient, error) {
	cfg, err := loadConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("✅ Connected to Redis successfully")
	return &redisClient{client: client}, nil
}

// Set wraps redis Set command.
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get wraps redis Get command.
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Ping wraps redis Ping command.
func (r *redisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
