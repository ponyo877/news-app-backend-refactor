package comment

import (
	"strings"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Service Comment usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateComment create a Comment
func (s *Service) CreateComment(commnet entity.Comment) (entity.ID, error) {
	return s.repo.Create(commnet)
}

// GetComment get a Comment
func (s *Service) GetComment(id entity.ID) (*entity.Comment, error) {
	comments, err := s.repo.Get(id)
	if comments == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// SearchComments search Comment
func (s *Service) SearchComments(query string) ([]*entity.Comment, error) {
	comments, err := s.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, entity.ErrNotFound
	}
	return comments, nil
}

// ListComments list Comment
func (s *Service) ListComments() ([]*entity.Comment, error) {
	comments, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, entity.ErrNotFound
	}
	return comments, nil
}

// DeleteComment Delete a Comment
func (s *Service) DeleteComment(id entity.ID) error {
	if _, err := s.GetComment(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// UpdateComment Update a Comment
func (s *Service) UpdateComment(e entity.Comment) error {
	if err := e.Validate(); err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
