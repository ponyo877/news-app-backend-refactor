package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
)

// stubArticleUseCase GetArticleだけ差し替えるテスト用スタブ。
// 他メソッドは未使用のため埋め込みnilインターフェースのままでよい
type stubArticleUseCase struct {
	article.UseCase
	article entity.Article
	err     error
}

func (s *stubArticleUseCase) GetArticle(entity.ID) (entity.Article, error) {
	return s.article, s.err
}

const testAPRoot = "https://matome.example.com"

func newTestArticle(t *testing.T, title string) entity.Article {
	t.Helper()
	imageURL, err := entity.NewImageURL("https://example.com/image.jpg")
	if err != nil {
		t.Fatalf("NewImageURLが失敗しました: %v", err)
	}
	return entity.Article{
		ID:       entity.NewID(),
		Title:    entity.NewArticleTitle(title),
		URL:      "https://example.com/article",
		ImageURL: imageURL,
		Site: entity.Site{
			ID:    entity.NewID(),
			Title: "テストサイト",
		},
		PublishedAt: time.Now(),
	}
}

func performLanding(service article.UseCase, path string) *httptest.ResponseRecorder {
	e := echo.New()
	MakeLandingHandlers(e, service, testAPRoot)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestShowArticleLanding(t *testing.T) {
	testArticle := newTestArticle(t, "今日の人気記事タイトル")
	rec := performLanding(&stubArticleUseCase{article: testArticle}, "/a/"+testArticle.ID.String())

	if rec.Code != http.StatusOK {
		t.Fatalf("ステータスが200ではありません: %d", rec.Code)
	}
	if got := rec.Header().Get("Cache-Control"); got != "public, max-age=0, s-maxage=86400" {
		t.Fatalf("Cache-Controlが想定と異なります: %s", got)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`og:title" content="今日の人気記事タイトル"`,
		`og:image" content="https://example.com/image.jpg"`,
		`og:url" content="` + testAPRoot + `/a/` + testArticle.ID.String() + `"`,
		`matomekun://a/` + testArticle.ID.String(),
		`app-id=1546424384`,
		`テストサイト`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("レスポンスに %q が含まれていません", want)
		}
	}
}

func TestShowArticleLandingEscapesTitle(t *testing.T) {
	testArticle := newTestArticle(t, `<script>alert("x")</script>`)
	rec := performLanding(&stubArticleUseCase{article: testArticle}, "/a/"+testArticle.ID.String())

	if strings.Contains(rec.Body.String(), `<script>alert`) {
		t.Fatal("タイトルがエスケープされていません(XSS)")
	}
}

func TestShowArticleLandingNotFound(t *testing.T) {
	rec := performLanding(&stubArticleUseCase{err: entity.ErrNotFound}, "/a/"+entity.NewID().String())

	if rec.Code != http.StatusNotFound {
		t.Fatalf("ステータスが404ではありません: %d", rec.Code)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("404のCache-Controlはno-storeであるべきです: %s", got)
	}
	if !strings.Contains(rec.Body.String(), "記事が見つかりませんでした") {
		t.Error("404ページの文言がありません")
	}
}

func TestShowArticleLandingInvalidID(t *testing.T) {
	rec := performLanding(&stubArticleUseCase{}, "/a/not-a-uuid")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("不正IDでステータスが404ではありません: %d", rec.Code)
	}
}

func TestGetArticleMeta(t *testing.T) {
	testArticle := newTestArticle(t, "メタ取得テスト")
	e := echo.New()
	MakeArticleHandlers(e, &stubArticleUseCase{article: testArticle}, func(next echo.HandlerFunc) echo.HandlerFunc { return next })

	req := httptest.NewRequest(http.MethodGet, "/v1/article/meta/"+testArticle.ID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("ステータスが200ではありません: %d", rec.Code)
	}
	for _, want := range []string{`"titles":"メタ取得テスト"`, `"id":"` + testArticle.ID.String() + `"`} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("レスポンスに %q が含まれていません", want)
		}
	}

	req404 := httptest.NewRequest(http.MethodGet, "/v1/article/meta/not-a-uuid", nil)
	rec404 := httptest.NewRecorder()
	e.ServeHTTP(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Fatalf("不正IDでステータスが404ではありません: %d", rec404.Code)
	}
}
