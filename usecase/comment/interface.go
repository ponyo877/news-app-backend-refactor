package comment

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.Comment, error)
	Search(query string) ([]*entity.Comment, error)
	List() ([]*entity.Comment, error)
}

// Writer interface
type Writer interface {
	Create(e entity.Comment) (entity.ID, error)
	Update(e entity.Comment) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	CreateComment(message string, site entity.Site) (entity.ID, error)
	GetComment(id entity.ID) (*entity.Comment, error)
	SearchComments(query string) ([]*entity.Comment, error)
	ListComments(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]*entity.Comment, error)
	DeleteComment(id entity.ID) error
	UpdateComment(e *entity.Comment) error
}
