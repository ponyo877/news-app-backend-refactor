package user

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.User, error)
	Search(query string) ([]*entity.User, error)
	List() ([]*entity.User, error)
}

// Writer interface
type Writer interface {
	Create(e entity.User) (entity.ID, error)
	Update(e entity.User) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	CreateUser(name, avatarURL, deviceHash string) (entity.ID, error)
	GetUser(id entity.ID) (*entity.User, error)
	SearchUsers(query string) ([]*entity.User, error)
	ListUsers(baseCreatedAt time.Time, invisibleIDSet entity.IDSet) ([]*entity.User, error)
	DeleteUser(id entity.ID) error
	UpdateUser(e *entity.User) error
}
