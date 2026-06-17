package main

import (
	"github.com/rs/zerolog/log"
	"lesson01-ebvn/internal/api"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/pkg/logger"
	"lesson01-ebvn/pkg/redis"
)

// @title bookMark API
// @version 1.0
// @description API document for bookMark API
// @host localhost:8080
// @BasePath /
func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("fail to load config")
	}
	// set log level
	logger.SetLogLevel(cfg.LogLevel)
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create redis client")

	}
	engine := api.NewEngine(cfg, redisClient)

	if err := engine.Start(); err != nil {
		panic(err)
	}
}
