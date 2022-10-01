package article

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Get(ID entity.ID) (entity.Article, error)
	// Search(keyword string) ([]entity.Article, error)
	SearchOnlyID(keyword string) ([]entity.ID, error)
	List(IDList []entity.ID) ([]entity.Article, error)
	ListOption(basePublishedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error)
	ListOnlyIDOrderByViewCount(period string) ([]entity.ID, error)
}

// Writer interface
type Writer interface {
	Create(e entity.Article) (entity.ID, time.Time, error)
	CreateForSearch(e entity.Article) error
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
	CreateArticle(entity.Article) error
	GetArticle(ID entity.ID) (entity.Article, error)
	SearchArticles(query string) ([]entity.Article, error)
	ListArticles(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error)
	ListPopularArticles(period string) ([]entity.Article, error)
	IncrementViewCount(ID entity.ID) error
}