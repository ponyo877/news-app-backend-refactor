package fileio

import (
	"net/url"
	"path"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Service WebDAV usecase
type Service struct {
	repo Repository
	root string
}

// NewService create new service
func NewService(r Repository, root string) *Service {
	return &Service{
		repo: r,
		root: root,
	}
}

func (s *Service) SaveImage(e entity.Image) (string, error) {
	filePath, err := s.repo.Upload(e)
	if err != nil {
		return "", err
	}
	URL, err := url.Parse(s.root)
	if err != nil {
		return "", err
	}
	URL.Path = path.Join(URL.Path, filePath)
	return URL.String(), nil
}

func (s *Service) FetchImage(name string) (entity.Image, error) {
	image, err := s.repo.Download(name)
	if err != nil {
		return entity.Image{}, err
	}
	return image, nil
}
