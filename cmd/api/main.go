package main

import (
	"lesson01-ebvn/internal/api"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/pkg/redis"
)

// @title bookMark API
// @version 1.0
// @description API document for bookMark API
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("config is nil")
	}
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		panic("redis client is nil")
	}
	engine := api.NewEngine(cfg, redisClient)

	if err := engine.Start(); err != nil {
		panic(err)
	}
}
