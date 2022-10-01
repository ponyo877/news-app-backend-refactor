package article

import (
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
func (s *Service) CreateArticle(article entity.Article) error {
	if _, _, err := s.repository.Create(article); err != nil {
		return err
	}
	if err := s.repository.CreateForSearch(article); err != nil {
		return err
	}
	return nil
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
func (s *Service) SearchArticles(keyword string) ([]entity.Article, error) {
	IDList, err := s.repository.SearchOnlyID(keyword)
	if err != nil {
		return nil, err
	}
	if len(IDList) == 0 {
		return nil, entity.ErrNotFound
	}
	articles, err := s.repository.List(IDList)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// ListArticles list article
func (s *Service) ListArticles(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error) {
	articles, err := s.repository.ListOption(baseCreatedAt, invisibleIDSet)
	if err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return nil, entity.ErrNotFound
	}
	return articles, nil
}

func (s *Service) ListPopularArticles(period string) ([]entity.Article, error) {
	IDList, err := s.repository.ListOnlyIDOrderByViewCount(period)
	if err != nil {
		return nil, err
	}
	if len(IDList) == 0 {
		return nil, entity.ErrNotFound
	}
	articles, err := s.repository.List(IDList)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (s *Service) IncrementViewCount(id entity.ID) error {
	if err := s.repository.IncrementViewCount(id); err != nil {
		return err
	}
	return nil
}
