package service

//go:generate mockery --name BookMarkService --filename book_mark_service.go
type BookMarkService interface {
	GetHealthInfo() (string, string)
}
