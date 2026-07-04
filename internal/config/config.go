package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration values parsed from environment variables.
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark-service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        string `envconfig:"PORT" default:"8080"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
	HostName    string `envconfig:"HOST_NAME" default:"localhost:8080"`
}

// LoadConfig reads configuration from env. If InstanceID is empty, it auto-generates a unique UUID.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process env: %w", err)
	}
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}
	return &cfg, nil
}
