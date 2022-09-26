package user

import (
	"strings"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Service User usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateUser create a User
func (s *Service) CreateUser(name, avatarURL, deviceHash string) (entity.ID, error) {
	user, err := entity.NewUser(name, avatarURL, deviceHash)
	if err != nil {
		return entity.NewID(), err
	}
	return s.repo.Create(user)
}

// GetUser get a User
func (s *Service) GetUser(id entity.ID) (*entity.User, error) {
	user, err := s.repo.Get(id)
	if user == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SearchUsers search User
func (s *Service) SearchUsers(query string) ([]*entity.User, error) {
	user, err := s.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, entity.ErrNotFound
	}
	return user, nil
}

// ListUsers list User
func (s *Service) ListUsers() ([]*entity.User, error) {
	user, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, entity.ErrNotFound
	}
	return user, nil
}

// DeleteUser Delete a User
func (s *Service) DeleteUser(id entity.ID) error {
	if _, err := s.GetUser(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// UpdateUser Update a User
func (s *Service) UpdateUser(e entity.User) error {
	if err := e.Validate(); err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
