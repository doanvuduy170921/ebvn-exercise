package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/internal/repository/mocks"
	"testing"
)

func TestGetHealthInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                string
		config              *config.Config
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name: "success with full config",
			config: &config.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "abc-123-xyz",
				Port:        "8080",
			},
			expectedServiceName: "bookmark-service",
			expectedInstanceID:  "abc-123-xyz",
		},
		{
			name: "empty_service_name",
			config: &config.Config{
				ServiceName: "",
				InstanceID:  "abc-123-xyz",
				Port:        "8080",
			},
			expectedServiceName: "",
			expectedInstanceID:  "abc-123-xyz",
		},
		{
			name: "empty_instance_id",
			config: &config.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "",
				Port:        "8080",
			},
			expectedServiceName: "bookmark-service",
			expectedInstanceID:  "",
		},
		{
			name:                "empty_config",
			config:              &config.Config{},
			expectedServiceName: "",
			expectedInstanceID:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := mocks.NewUrlRepo(t)
			svc := NewBookMarkService(tc.config, mockRepo)
			// Act
			serviceName, instanceID := svc.GetHealthInfo()
			// Assert
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}
}

func TestGenerateKey(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		ServiceName: "bookmark-service",
		InstanceID:  "abc-123-xyz",
		Port:        "8080",
	}

	testCases := []struct {
		name        string
		url         string
		exp         int
		setupMock   func(ctx context.Context) *mocks.UrlRepo
		expectError bool
	}{
		{
			name: "success",
			url:  "https://google.com",
			exp:  604800,
			setupMock: func(ctx context.Context) *mocks.UrlRepo {
				m := mocks.NewUrlRepo(t)
				m.On("Save", ctx, mock.AnythingOfType("string"), "https://google.com", 604800).
					Return(nil)
				return m
			},
			expectError: false,
		},
		{
			name: "save error",
			url:  "https://google.com",
			exp:  604800,
			setupMock: func(ctx context.Context) *mocks.UrlRepo {
				m := mocks.NewUrlRepo(t)
				m.On("Save", ctx, mock.AnythingOfType("string"), "https://google.com", 604800).
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
			mockRepo := tc.setupMock(ctx)
			svc := NewBookMarkService(cfg, mockRepo)

			code, err := svc.GenerateKey(ctx, tc.url, tc.exp)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.Len(t, code, 7)
			}
		})
	}
}
