package repository

import (
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
)

// CommentMySQL mysql repository
type CommentMySQL struct {
	db *gorm.DB
}

type CommentMySQLPresenter struct {
	ID         string    `gorm:"column:id;primary_key"`
	ArticleID  string    `gorm:"column:article_id"`
	UserName   string    `gorm:"column:user_name"`
	ImageURL   string    `gorm:"column:image_url"`
	DeviceHash string    `gorm:"column:device_hash"`
	Message    string    `gorm:"column:message"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

// TableName
func (s CommentMySQLPresenter) TableName() string {
	return "comments"
}

// NewCommentMySQL create new repository
func NewCommentMySQL(db *gorm.DB) *CommentMySQL {
	return &CommentMySQL{
		db: db,
	}
}

// Get
func (r *CommentMySQL) Get(id entity.ID) (entity.Comment, error) {
	return entity.Comment{}, nil
}

// Search
func (r *CommentMySQL) Search(query string) ([]entity.Comment, error) {
	return []entity.Comment{}, nil
}

// List
func (r *CommentMySQL) List(articleID entity.ID) ([]entity.Comment, error) {
	var commentMySQLList []CommentMySQLPresenter
	err := r.db.Where("article_id = ?", articleID.String()).Find(&commentMySQLList).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Infof("DBの接続に失敗しました: %v", err)
		return nil, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return []entity.Comment{}, nil
	}
	var commentList []entity.Comment
	for _, commentMySQL := range commentMySQLList {
		ID, err := entity.StringToID(commentMySQL.ID)
		if err != nil {
			return nil, err
		}
		ArticleID, err := entity.StringToID(commentMySQL.ArticleID)
		if err != nil {
			return nil, err
		}
		comment := entity.Comment{
			ID:         ID,
			UserName:   commentMySQL.UserName,
			AvatarURL:  commentMySQL.ImageURL,
			DeviceHash: commentMySQL.DeviceHash,
			Message:    commentMySQL.Message,
			Article: entity.Article{
				ID: ArticleID,
			},
			UpdatedAt: commentMySQL.UpdatedAt,
			CreatedAt: commentMySQL.CreatedAt,
		}
		commentList = append(commentList, comment)
	}
	return commentList, nil
}

// Create
func (r *CommentMySQL) Create(e entity.Comment) (entity.ID, error) {
	commentMySQLPresenter := CommentMySQLPresenter{
		ID:         e.ID.String(),
		ArticleID:  e.Article.ID.String(),
		UserName:   e.UserName,
		ImageURL:   e.AvatarURL,
		DeviceHash: e.DeviceHash,
		Message:    e.Message,
		UpdatedAt:  e.UpdatedAt,
		CreatedAt:  e.CreatedAt,
	}
	if err := r.db.Create(commentMySQLPresenter).Error; err != nil {
		return entity.NewID(), nil
	}
	return e.ID, nil
}

// Update
func (r *CommentMySQL) Update(e entity.Comment) error {
	return nil
}

// Delete
func (r *CommentMySQL) Delete(id entity.ID) error {
	return nil
}
