package site

import (
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Get(id entity.ID) (entity.Site, error)
	Search(query string) ([]entity.Site, error)
	List() ([]entity.Site, error)
}

// Writer interface
type Writer interface {
	Create(e entity.Site) (entity.ID, error)
	Update(e entity.Site) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	CreateSite(title string, RSSURL string, ImageURL string) (entity.ID, error)
	GetSite(id entity.ID) (entity.Site, error)
	ListSite() ([]entity.Site, error)
	UpdateSite(e entity.Site) error
}
