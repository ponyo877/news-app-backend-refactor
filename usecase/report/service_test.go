package report

import (
	"errors"
	"testing"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type fakeRepository struct {
	saved []entity.ArticleReport
	err   error
}

func (f *fakeRepository) Save(e entity.ArticleReport) error {
	if f.err != nil {
		return f.err
	}
	f.saved = append(f.saved, e)
	return nil
}

func TestReportArticle(t *testing.T) {
	articleID := entity.NewID()

	t.Run("正常な報告は1件保存される", func(t *testing.T) {
		repo := &fakeRepository{}
		service := NewService(repo)
		err := service.ReportArticle(articleID, "https://example.com/1", "テスト速報", "hash-1", "android", "1.52", 4, "ad_remains")
		if err != nil {
			t.Fatalf("エラーになりました: %v", err)
		}
		if len(repo.saved) != 1 {
			t.Fatalf("保存件数が1ではありません: %d", len(repo.saved))
		}
		saved := repo.saved[0]
		if saved.ArticleID.String() != articleID.String() || saved.Reason != "ad_remains" || saved.RulesVersion != 4 || saved.DeviceHash != "hash-1" {
			t.Fatalf("保存内容が一致しません: %+v", saved)
		}
	})

	t.Run("理由が不正なら保存せずErrInvalidEntity", func(t *testing.T) {
		repo := &fakeRepository{}
		service := NewService(repo)
		err := service.ReportArticle(articleID, "https://example.com/1", "テスト速報", "hash-1", "ios", "1.52", 4, "spam")
		if !errors.Is(err, entity.ErrInvalidEntity) {
			t.Fatalf("ErrInvalidEntityが返っていません: %v", err)
		}
		if len(repo.saved) != 0 {
			t.Fatalf("不正な報告が保存されました")
		}
	})

	t.Run("保存の失敗はそのまま返る", func(t *testing.T) {
		dbErr := errors.New("db down")
		service := NewService(&fakeRepository{err: dbErr})
		err := service.ReportArticle(articleID, "https://example.com/1", "テスト速報", "hash-1", "ios", "1.52", 4, "other")
		if !errors.Is(err, dbErr) {
			t.Fatalf("保存エラーが伝播していません: %v", err)
		}
	})
}
