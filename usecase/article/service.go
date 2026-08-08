package article

import (
	"time"

	"github.com/labstack/gommon/log"
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
	return nil
}

// GetArticle get a article。
// 以前はErrNotFoundをnil+ゼロ値に握り潰していたが、着地ページ(/a/)とメタAPIが
// 404を判定できるようそのまま返す(修正時点で他に呼び出し元なし)
func (s *Service) GetArticle(id entity.ID) (entity.Article, error) {
	return s.repository.Get(id)
}

// SearchArticles search article
func (s *Service) SearchArticles(keyword entity.Keyword) ([]entity.Article, error) {
	return s.repository.SearchOnlyID(keyword)
}

// ListArticles list article
func (s *Service) ListArticles(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error) {
	return s.repository.ListOption(baseCreatedAt, invisibleIDSet, 15)
}

// ListPopularArticles
func (s *Service) ListPopularArticles(period string) ([]entity.Article, error) {
	IDList, err := s.repository.ListOnlyIDOrderByViewCount(period)
	if err != nil {
		return nil, err
	}
	if len(IDList) == 0 {
		return []entity.Article{}, nil
	}
	return s.repository.List(IDList)
}

// IncrementViewCount
func (s *Service) IncrementViewCount(id entity.ID) error {
	return s.repository.IncrementViewCount(id)
}

func (s *Service) ListSimilarArticles(ID entity.ID) ([]entity.Article, error) {
	IDList, err := s.repository.ListBySimilarity(ID)
	if err == entity.ErrNotFound {
		return []entity.Article{}, nil
	}
	if err != nil {
		return []entity.Article{}, err
	}
	if len(IDList) == 0 {
		return []entity.Article{}, nil
	}
	return s.repository.List(IDList)
}

func (s *Service) UpdateMLIndex() error {
	targetArticles, err := s.repository.ListOption(time.Time{}, entity.NewIDSet(), 100)
	if err != nil {
		return err
	}
	log.Infof("UpdateMLIndexで %v 個の記事がMLIndexに登録されます", len(targetArticles))
	return s.repository.CreateMLIndex(targetArticles)
}
