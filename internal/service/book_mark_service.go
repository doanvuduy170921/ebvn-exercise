package service

import (
	"context"
	"fmt"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/internal/repository"
	"lesson01-ebvn/utils"
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
