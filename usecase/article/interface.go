package article

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Get(ID entity.ID) (entity.Article, error)
	Search(query string) ([]entity.Article, error)
	List(basePublishedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error)
	ListOrderByViewCount(period string) ([]entity.Article, error)
}

// Writer interface
type Writer interface {
	Create(e entity.Article) (entity.ID, time.Time, error)
	Update(e entity.Article) error
	DeleteByID(ID entity.ID) error
	IncrementViewCount(ID entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	CreateArticle(entity.Article) (entity.ID, time.Time, error)
	GetArticle(ID entity.ID) (entity.Article, error)
	SearchArticles(query string) ([]entity.Article, error)
	ListArticles(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error)
	ListPopularArticles(period string) ([]entity.Article, error)
	IncrementViewCount(ID entity.ID) error
}
