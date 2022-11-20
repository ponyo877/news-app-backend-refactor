package article

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Get(ID entity.ID) (entity.Article, error)
	SearchOnlyID(keyword entity.Keyword) ([]entity.Article, error)
	// SearchOnlyID(keyword string) ([]entity.ID, error) // ElasticSearch用
	List(IDList []entity.ID) ([]entity.Article, error)
	ListOption(basePublishedAt time.Time, invisibleIDSet entity.IDSet, limit int) ([]entity.Article, error)
	ListOnlyIDOrderByViewCount(period string) ([]entity.ID, error)
	ListBySimilarity(ID entity.ID) ([]entity.ID, error)
	// GetArticleNumberByArticleID(articleID entity.ID, prefix string) (int, error)
}

// Writer interface
type Writer interface {
	Create(e entity.Article) (entity.ID, time.Time, error)
	// CreateForSearch(e entity.Article) error // ElasticSearch用
	Update(e entity.Article) error
	DeleteByID(ID entity.ID) error
	IncrementViewCount(ID entity.ID) error
	CreateMLIndex(articles []entity.Article) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	CreateArticle(entity.Article) error
	GetArticle(ID entity.ID) (entity.Article, error)
	SearchArticles(keyword entity.Keyword) ([]entity.Article, error)
	ListArticles(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error)
	ListPopularArticles(period string) ([]entity.Article, error)
	IncrementViewCount(ID entity.ID) error
	ListSimilarArticles(ID entity.ID) ([]entity.Article, error)
	UpdateMLIndex() error
}
