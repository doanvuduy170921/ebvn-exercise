package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/internal/repository"
	"lesson01-ebvn/utils"
)

var (
	ErrorNotFound = errors.New("not found")
)

// bookMarkService implements BookMarkService with config and repository dependencies.
type bookMarkService struct {
	cfg  *config.Config
	repo repository.UrlRepo
}

// NewBookMarkService creates a new bookMarkService with the given config and repository.
func NewBookMarkService(cfg *config.Config, repo repository.UrlRepo) BookMarkService {
	return &bookMarkService{
		cfg:  cfg,
		repo: repo,
	}

}

// GetHealthInfo returns the service name and instance ID from config.
func (b *bookMarkService) GetHealthInfo() (string, string) {
	return b.cfg.ServiceName, b.cfg.InstanceID
}

// GenerateKey creates a random 7-character short code, saves the code-URL pair to Redis, and returns the code.
func (b *bookMarkService) GenerateKey(ctx context.Context, url string, exp int) (string, error) {
	code, err := utils.GenerateShortCode(7)
	if err != nil {
		return "", fmt.Errorf("generate key code error: %v", err)
	}

	err = b.repo.Save(ctx, code, url, exp)
	if err != nil {
		return "", fmt.Errorf("save bookmark url error: %v", err)
	}
	return code, nil

}

func (b *bookMarkService) GetURL(ctx context.Context, code string) (string, error) {
	url, err := b.repo.Get(ctx, code)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrorNotFound
		}
	}
	return url, nil
}
