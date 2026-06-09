package config

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config holds all application configuration values parsed from environment variables.
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark-service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        string `envconfig:"PORT" default:"8080"`
}

// LoadConfig reads configuration from env. If InstanceID is empty, it auto-generates a unique UUID.
func LoadConfig() *Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}
	return &cfg
}
