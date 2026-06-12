package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client using environment-based configuration.
// It pings Redis to verify the connection before returning.
func NewRedisClient() (*redis.Client, error) {
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
	return client, nil
}
