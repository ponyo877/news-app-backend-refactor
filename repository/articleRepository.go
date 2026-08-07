package repository

import (
	"github.com/go-redis/redis/v9"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
	"gorm.io/gorm"
)

// ArticleRepository mysql repository
type ArticleRepository struct {
	db        *gorm.DB
	kvs       *redis.Client
	model     textencoding.Interface
	index     VectorIndex
	indexPath string
}

// NewArticleRepository create new repository
func NewArticleRepository(db *gorm.DB, kvs *redis.Client, model textencoding.Interface, index VectorIndex, indexPath string) *ArticleRepository {
	return &ArticleRepository{
		db,
		kvs,
		model,
		index,
		indexPath,
	}
}
