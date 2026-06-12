package repository

import (
	"context"
	"errors"
	redis "lesson01-ebvn/pkg/redis/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	redisMock "lesson01-ebvn/pkg/redis/mocks"
)

func TestUrlRepo_Save(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		key         string
		value       string
		exp         int
		setupMock   func(ctx context.Context) *redisMock.RedisClient // ← sửa
		expectError bool
	}{
		{
			name:  "success",
			key:   "abc1234",
			value: "https://google.com",
			exp:   604800,
			setupMock: func(ctx context.Context) *redisMock.RedisClient {
				m := redis.NewRedisClient(t) // ← sửa
				m.On("Set", ctx, "abc1234", "https://google.com", time.Duration(604800)*time.Second).
					Return(nil)
				return m
			},
			expectError: false,
		},
		{
			name:  "redis error",
			key:   "abc1234",
			value: "https://google.com",
			exp:   604800,
			setupMock: func(ctx context.Context) *redisMock.RedisClient {
				m := redis.NewRedisClient(t) // ← sửa
				m.On("Set", ctx, "abc1234", "https://google.com", time.Duration(604800)*time.Second).
					Return(errors.New("redis error"))
				return m
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockRedis := tc.setupMock(ctx)
			repo := NewUrlRepo(mockRedis)

			err := repo.Save(ctx, tc.key, tc.value, tc.exp)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
