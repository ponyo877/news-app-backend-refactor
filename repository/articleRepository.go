package repository

import (
	"github.com/go-redis/redis/v9"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

// ArticleRepository mysql repository
type ArticleRepository struct {
	db  *gorm.DB
	kvs *redis.Client
	se  *elastic.Client
}

// NewArticleRepository create new repository
func NewArticleRepository(db *gorm.DB, kvs *redis.Client, se *elastic.Client) *ArticleRepository {
	return &ArticleRepository{
		db,
		kvs,
		se,
	}
}
