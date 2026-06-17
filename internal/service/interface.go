package service

import (
	"context"
)

//go:generate mockery --name BookMarkService --filename book_mark_service.go
type BookMarkService interface {
	GetHealthInfo() (string, string)
	GenerateKey(ctx context.Context, url string, exp int) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}
