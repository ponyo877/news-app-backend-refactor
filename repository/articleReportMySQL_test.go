package repository

import (
	"strings"
	"testing"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// upsertが対応状態の列を触らないことをSQL文で固定する(DBには接続しない)
func TestArticleReportMySQLSaveSQL(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{SkipInitializeWithVersion: true}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		// 既定ではCreate前にトランザクションを開くため、DryRunでもDBへ接続しに行ってしまう
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("gorm.Openが失敗しました: %v", err)
	}
	report := entity.ArticleReport{
		ArticleID:  entity.NewID(),
		URL:        "https://example.com/archives/1.html",
		SiteTitle:  "テスト速報",
		DeviceHash: "abc123",
		Platform:   "ios",
		AppVersion: "1.52",
		Reason:     "missing_media",
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
	presenter := ArticleReportMySQLPresenter{
		ArticleID:  report.ArticleID.String(),
		URL:        report.URL,
		SiteTitle:  report.SiteTitle,
		DeviceHash: report.DeviceHash,
		Platform:   report.Platform,
		AppVersion: report.AppVersion,
		Reason:     report.Reason,
		UpdatedAt:  report.UpdatedAt,
		CreatedAt:  report.CreatedAt,
	}
	stmt := db.Session(&gorm.Session{DryRun: true}).Clauses(upsertClause()).Create(&presenter).Statement
	sql := stmt.SQL.String()
	if !strings.Contains(sql, "ON DUPLICATE KEY UPDATE") {
		t.Fatalf("upsertになっていません: %s", sql)
	}
	for _, column := range []string{"status", "admin_note", "created_at`=VALUES"} {
		if strings.Contains(sql, column) {
			t.Fatalf("再報告で触ってはいけない列が含まれています(%s): %s", column, sql)
		}
	}
	for _, column := range []string{"`reason`=VALUES(`reason`)", "`updated_at`=VALUES(`updated_at`)"} {
		if !strings.Contains(sql, column) {
			t.Fatalf("再報告で上書きすべき列が含まれていません(%s): %s", column, sql)
		}
	}
}
