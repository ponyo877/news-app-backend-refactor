package site

import (
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Service Article usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateSite create a article
func (s *Service) CreateSite(title string, RSSURL string, ImageURL string) (entity.ID, error) {
	site, err := entity.NewSite(title, RSSURL, ImageURL)
	if err != nil {
		return entity.NewID(), err
	}
	return s.repo.Create(site)
}

// GetSite get a article
func (s *Service) GetSite(id entity.ID) (*entity.Site, error) {
	a, err := s.repo.Get(id)
	if a == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ListSite list article
func (s *Service) ListSite() ([]entity.Site, error) {
	sites, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(sites) == 0 {
		return nil, entity.ErrNotFound
	}
	return sites, nil
}

// UpdateSite Update a article
func (s *Service) UpdateSite(e *entity.Site) error {
	if err := e.Validate(); err != nil {
		return err
	}
	return s.repo.Update(e)
}
