package entity

import (
	"errors"
	"strings"
	"testing"
)

func TestNewArticleReport(t *testing.T) {
	articleID := NewID()
	valid := func() (ArticleReport, error) {
		return NewArticleReport(articleID, "https://example.com/archives/1.html", "テスト速報", "abc123", "ios", "1.52", 4, "missing_media")
	}

	report, err := valid()
	if err != nil {
		t.Fatalf("正常な入力でエラーになりました: %v", err)
	}
	if report.Reason != "missing_media" || report.RulesVersion != 4 || report.SiteTitle != "テスト速報" {
		t.Fatalf("フィールドが保持されていません: %+v", report)
	}
	if report.CreatedAt.IsZero() || report.UpdatedAt.IsZero() {
		t.Fatalf("日時が設定されていません: %+v", report)
	}

	// サイト名は拒否せず100文字に切り詰める(マルチバイト単位)
	long, err := NewArticleReport(articleID, "https://example.com/", strings.Repeat("あ", 150), "abc123", "android", "", 0, "other")
	if err != nil {
		t.Fatalf("長いサイト名が拒否されました: %v", err)
	}
	if got := len([]rune(long.SiteTitle)); got != 100 {
		t.Fatalf("サイト名が切り詰められていません: %d文字", got)
	}

	cases := map[string]struct {
		url, deviceHash, platform, appVersion, reason string
		rulesVersion                                  int
	}{
		"理由が未知":          {"https://example.com/", "abc123", "ios", "1.52", "spam", 4},
		"理由が空":           {"https://example.com/", "abc123", "ios", "1.52", "", 4},
		"URLが空":          {"", "abc123", "ios", "1.52", "other", 4},
		"URLがhttp以外":     {"javascript:alert(1)", "abc123", "ios", "1.52", "other", 4},
		"URLが長すぎる":       {"https://example.com/" + strings.Repeat("a", 2048), "abc123", "ios", "1.52", "other", 4},
		"devicehashが空":   {"https://example.com/", "", "ios", "1.52", "other", 4},
		"devicehashが長い":  {"https://example.com/", strings.Repeat("a", 65), "ios", "1.52", "other", 4},
		"platformが不正":    {"https://example.com/", "abc123", "web", "1.52", "other", 4},
		"appversionが長い":  {"https://example.com/", "abc123", "ios", strings.Repeat("1", 21), "other", 4},
		"rulesversionが負": {"https://example.com/", "abc123", "ios", "1.52", "other", -1},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := NewArticleReport(articleID, c.url, "site", c.deviceHash, c.platform, c.appVersion, c.rulesVersion, c.reason)
			if !errors.Is(err, ErrInvalidEntity) {
				t.Fatalf("ErrInvalidEntityが返っていません: %v", err)
			}
		})
	}
}
