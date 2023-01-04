package comment

import (
	"strings"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/user"
)

// Service Comment usecase
type Service struct {
	repository  Repository
	userService user.UseCase
}

// NewService create new service
func NewService(r Repository, u user.UseCase) *Service {
	return &Service{
		repository:  r,
		userService: u,
	}
}

// CreateComment create a Comment
func (s *Service) CreateComment(commnet entity.Comment) (entity.ID, error) {
	user, err := s.userService.GetUserOption(commnet.DeviceHash)
	if err != nil {
		return entity.NewID(), err
	}
	commnet.UserName = user.Name
	commnet.AvatarURL = user.AvatarURL
	return s.repository.Create(commnet)
}

// GetComment get a Comment
func (s *Service) GetComment(id entity.ID) (entity.Comment, error) {
	comments, err := s.repository.Get(id)
	if err == entity.ErrNotFound {
		return comments, nil
	}
	if err != nil {
		return entity.Comment{}, err
	}
	return comments, nil
}

// SearchComments search Comment
func (s *Service) SearchComments(query string) ([]entity.Comment, error) {
	return s.repository.Search(strings.ToLower(query))
}

// ListComments list Comment
func (s *Service) ListComments(articleID entity.ID) ([]entity.Comment, error) {
	return s.repository.List(articleID)
}

// DeleteComment Delete a Comment
func (s *Service) DeleteComment(id entity.ID) error {
	if _, err := s.GetComment(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}

// UpdateComment Update a Comment
func (s *Service) UpdateComment(e entity.Comment) error {
	if err := e.Validate(); err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repository.Update(e)
}
