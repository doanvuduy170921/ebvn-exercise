package service

import "lesson01-ebvn/internal/config"

type bookMarkService struct {
	cfg *config.Config
}

func NewBookMarkService(cfg *config.Config) BookMarkService {
	return &bookMarkService{
		cfg: cfg,
	}

}

func (b *bookMarkService) GetHealthInfo() (string, string) {
	return b.cfg.ServiceName, b.cfg.InstanceID
}
