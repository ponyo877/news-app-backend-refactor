package article

import (
	"strings"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Service Article usecase
type Service struct {
	repository Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

// CreateArticle create a article
func (s *Service) CreateArticle(article entity.Article) (entity.ID, time.Time, error) {
	return s.repository.Create(article)
}

// GetArticle get a article
func (s *Service) GetArticle(id entity.ID) (entity.Article, error) {
	article, err := s.repository.Get(id)
	if err != nil {
		return entity.Article{}, err
	}
	return article, nil
}

// SearchArticles search article
func (s *Service) SearchArticles(query string) ([]entity.Article, error) {
	articles, err := s.repository.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return nil, entity.ErrNotFound
	}
	return articles, nil
}

// ListArticles list article
func (s *Service) ListArticles(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error) {
	articles, err := s.repository.List(baseCreatedAt, invisibleIDSet)
	if err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return nil, entity.ErrNotFound
	}
	return articles, nil
}

func (s *Service) ListPopularArticles(period string) ([]entity.Article, error) {
	articles, err := s.repository.ListOrderByViewCount(period)
	if err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return nil, entity.ErrNotFound
	}
	return articles, nil
}

func (s *Service) IncrementViewCount(id entity.ID) error {
	if err := s.repository.IncrementViewCount(id); err != nil {
		return err
	}
	return nil
}
