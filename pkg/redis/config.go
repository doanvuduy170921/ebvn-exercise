package redis

import "github.com/kelseyhightower/envconfig"

// Config holds the configuration for connecting to Redis.
type Config struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

// loadConfig loads Redis configuration from environment variables with the given prefix.
func loadConfig(preFix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(preFix, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
