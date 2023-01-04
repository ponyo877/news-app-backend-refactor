package fileio

import (
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Download(name string) (entity.Image, error)
}

// Writer interface
type Writer interface {
	Upload(e entity.Image) (string, error)
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	SaveImage(entity.Image) (string, error)
	FetchImage(name string) (entity.Image, error)
}
