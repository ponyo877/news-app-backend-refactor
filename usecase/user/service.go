package user

import (
	"strings"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/fileio"
)

// Service User usecase
type Service struct {
	repository    Repository
	fileioService fileio.UseCase
}

// NewService create new service
func NewService(r Repository, f fileio.UseCase) *Service {
	return &Service{
		repository:    r,
		fileioService: f,
	}
}

// CreateUser create a User
func (s *Service) CreateUser(name string, avatarImage entity.Image, deviceHash string) (entity.ID, error) {
	avatarURL, err := s.fileioService.SaveImage(avatarImage)
	if err != nil {
		return entity.NewID(), err
	}
	user, err := entity.NewUser(name, avatarURL, deviceHash)
	if err != nil {
		return entity.NewID(), err
	}
	return s.repository.Create(user)
}

// GetUser get a User
func (s *Service) GetUser(ID entity.ID) (entity.User, error) {
	user, err := s.repository.Get(ID)
	if err == entity.ErrNotFound {
		return entity.User{}, nil
	}
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// GetUser get a User
func (s *Service) GetUserOption(deviceHash string) (entity.User, error) {
	user, err := s.repository.GetOption(deviceHash)
	if err == entity.ErrNotFound {
		return entity.User{}, nil
	}
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// SearchUsers search User
func (s *Service) SearchUsers(query string) ([]entity.User, error) {
	user, err := s.repository.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers list User
func (s *Service) ListUsers() ([]entity.User, error) {
	user, err := s.repository.List()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser Delete a User
func (s *Service) DeleteUser(id entity.ID) error {
	if _, err := s.GetUser(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}

// UpdateUser Update a User
func (s *Service) UpdateUser(e entity.User) error {
	if err := e.Validate(); err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repository.Update(e)
}
