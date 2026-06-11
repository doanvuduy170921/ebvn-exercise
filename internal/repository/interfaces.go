package repository

import "context"

//go:generate mockery --name UrlRepo --filename=url_repo.go --outpkg mocks
type UrlRepo interface {
	Save(ctx context.Context, key, value string, exp int) error
	Get(ctx context.Context, key string) (string, error)
}
