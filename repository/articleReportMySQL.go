package repository

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ArticleReportMySQL mysql repository
type ArticleReportMySQL struct {
	db *gorm.DB
}

// ArticleReportMySQLPresenter article_reports テーブルのうちアプリが書く列。
//
// 対応状態の列(status / admin_note / status_updated_at)は意図的に持たない。
// 管理者がSQLで更新する列で、INSERT時はDBの DEFAULT 'open' に任せ、
// 再報告のupsertでも構造的に触れないようにするため(再報告で対応済みが未対応に戻らない)
type ArticleReportMySQLPresenter struct {
	ID           uint64    `gorm:"column:id;primary_key;auto_increment"`
	ArticleID    string    `gorm:"column:article_id"`
	URL          string    `gorm:"column:url"`
	SiteTitle    string    `gorm:"column:site_title"`
	DeviceHash   string    `gorm:"column:device_hash"`
	Platform     string    `gorm:"column:platform"`
	AppVersion   string    `gorm:"column:app_version"`
	RulesVersion int       `gorm:"column:rules_version"`
	Reason       string    `gorm:"column:reason"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

// TableName
func (s ArticleReportMySQLPresenter) TableName() string {
	return "article_reports"
}

// NewArticleReportMySQL create new repository
func NewArticleReportMySQL(db *gorm.DB) *ArticleReportMySQL {
	return &ArticleReportMySQL{
		db: db,
	}
}

// 同一端末×記事の再報告で上書きする列。created_at と対応状態の列は含めない
var articleReportUpsertColumns = []string{"url", "site_title", "platform", "app_version", "rules_version", "reason", "updated_at"}

// Save 報告を登録する。同一端末×同一記事(UNIQUE uq_device_article)は理由・バージョン・日時の上書き
func (r *ArticleReportMySQL) Save(e entity.ArticleReport) error {
	presenter := ArticleReportMySQLPresenter{
		ArticleID:    e.ArticleID.String(),
		URL:          e.URL,
		SiteTitle:    e.SiteTitle,
		DeviceHash:   e.DeviceHash,
		Platform:     e.Platform,
		AppVersion:   e.AppVersion,
		RulesVersion: e.RulesVersion,
		Reason:       e.Reason,
		UpdatedAt:    e.UpdatedAt,
		CreatedAt:    e.CreatedAt,
	}
	return r.db.Clauses(upsertClause()).Create(&presenter).Error
}

// upsertClause 同一端末×記事の再報告を上書きにする句。
// MySQLでは ON DUPLICATE KEY UPDATE になる(Columnsは無視されるがユニークキーは1本なので問題ない)
func upsertClause() clause.OnConflict {
	return clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_hash"}, {Name: "article_id"}},
		DoUpdates: clause.AssignmentColumns(articleReportUpsertColumns),
	}
}
