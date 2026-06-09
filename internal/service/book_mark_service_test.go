package service

import (
	"github.com/stretchr/testify/assert"
	"lesson01-ebvn/internal/config"
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

			svc := NewBookMarkService(tc.config)
			// Act
			serviceName, instanceID := svc.GetHealthInfo()
			// Assert
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}
}
