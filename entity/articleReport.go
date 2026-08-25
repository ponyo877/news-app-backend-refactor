package entity

import (
	"strings"
	"time"
)

// ArticleReport 記事表示の不具合報告(アプリの記事画面⋮メニュー「表示の不具合を報告」)。
//
// 1.51で記事内の埋め込み(Xポスト・imgur等)が消えていた際、ユーザーから具体的な記事URLが届かず
// 実態把握に時間がかかった。理由の選択だけで記事URL・アプリ内ID・端末情報を蓄積し、
// 管理者がSQLで閲覧して対応状態(status)を管理する。status系のカラムはDB側にだけ存在し、
// アプリからは一切触れない(repository/articleReportMySQL.go 参照)
type ArticleReport struct {
	ArticleID    ID
	URL          string
	SiteTitle    string
	DeviceHash   string
	Platform     string
	AppVersion   string
	RulesVersion int
	Reason       string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// ArticleReportReasons 報告理由。アプリ側 src/lib/articleReport.ts の key と一致させる
var ArticleReportReasons = []string{"missing_media", "ad_remains", "body_broken", "other"}

const (
	articleReportMaxURLLength        = 2048
	articleReportMaxSiteTitleLength  = 100
	articleReportMaxDeviceHashLength = 64
	articleReportMaxAppVersionLength = 20
)

// NewArticleReport create a new article report
func NewArticleReport(articleID ID, url, siteTitle, deviceHash, platform, appVersion string, rulesVersion int, reason string) (ArticleReport, error) {
	report := ArticleReport{
		ArticleID: articleID,
		URL:       url,
		// サイト名は集計用の付帯情報なので、長すぎても拒否せず切り詰める
		SiteTitle:    truncateSiteTitle(siteTitle),
		DeviceHash:   deviceHash,
		Platform:     platform,
		AppVersion:   appVersion,
		RulesVersion: rulesVersion,
		Reason:       reason,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}
	if err := report.Validate(); err != nil {
		return ArticleReport{}, ErrInvalidEntity
	}
	return report, nil
}

// Validate validate data
func (r *ArticleReport) Validate() error {
	if !isArticleReportReason(r.Reason) {
		return ErrInvalidEntity
	}
	// URLは管理者が開いて確認する値なので、http(s)以外(javascript:等の投げ込み)は入れない
	if r.URL == "" || len(r.URL) > articleReportMaxURLLength ||
		!(strings.HasPrefix(r.URL, "http://") || strings.HasPrefix(r.URL, "https://")) {
		return ErrInvalidEntity
	}
	if r.DeviceHash == "" || len(r.DeviceHash) > articleReportMaxDeviceHashLength {
		return ErrInvalidEntity
	}
	if r.Platform != "ios" && r.Platform != "android" {
		return ErrInvalidEntity
	}
	// 開発ビルドではアプリバージョンが取れないことがあるため空は許容する
	if len(r.AppVersion) > articleReportMaxAppVersionLength {
		return ErrInvalidEntity
	}
	if r.RulesVersion < 0 {
		return ErrInvalidEntity
	}
	return nil
}

func isArticleReportReason(reason string) bool {
	for _, candidate := range ArticleReportReasons {
		if reason == candidate {
			return true
		}
	}
	return false
}

func truncateSiteTitle(siteTitle string) string {
	runes := []rune(siteTitle)
	if len(runes) <= articleReportMaxSiteTitleLength {
		return siteTitle
	}
	return string(runes[:articleReportMaxSiteTitleLength])
}
