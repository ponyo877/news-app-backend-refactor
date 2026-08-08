package repository

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
)

// DigestLogMySQL mysql repository
type DigestLogMySQL struct {
	db *gorm.DB
}

type DigestLogMySQLPresenter struct {
	ID        uint64    `gorm:"column:id;primary_key;auto_increment"`
	ArticleID string    `gorm:"column:article_id"`
	SentCount int       `gorm:"column:sent_count"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// TableName
func (s DigestLogMySQLPresenter) TableName() string {
	return "digest_logs"
}

// NewDigestLogMySQL create new repository
func NewDigestLogMySQL(db *gorm.DB) *DigestLogMySQL {
	return &DigestLogMySQL{
		db: db,
	}
}

// Save 送信結果を記録する(重複送信防止とCTR計測の分母)
func (r *DigestLogMySQL) Save(articleID entity.ID, sentCount int) error {
	presenter := DigestLogMySQLPresenter{
		ArticleID: articleID.String(),
		SentCount: sentCount,
		CreatedAt: time.Now(),
	}
	return r.db.Create(&presenter).Error
}

// ListRecentArticleIDs 指定時刻以降に送信した記事IDの一覧
func (r *DigestLogMySQL) ListRecentArticleIDs(since time.Time) ([]string, error) {
	var articleIDs []string
	if err := r.db.
		Model(&DigestLogMySQLPresenter{}).
		Where("created_at > ?", since).
		Pluck("article_id", &articleIDs).
		Error; err != nil {
		return nil, err
	}
	return articleIDs, nil
}
