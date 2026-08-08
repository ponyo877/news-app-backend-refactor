package handler

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
)

// 共有着地ページ(/a/{id})。
// OGPクローラー(X/LINE/Discord)はJSを実行しないためサーバサイドレンダリング必須。
// テンプレートはembedでバイナリ同梱のため、CGO_ENABLED=0のクロスコンパイルに影響しない

//go:embed landing_article.html
var landingArticleHTML string

//go:embed landing_notfound.html
var landingNotFoundHTML string

var landingArticleTmpl = template.Must(template.New("landing").Parse(landingArticleHTML))

type landingView struct {
	ID           string
	Title        string
	SiteTitle    string
	ImageURL     string
	CanonicalURL string
}

// MakeLandingHandlers
func MakeLandingHandlers(e *echo.Echo, service article.UseCase, apRoot string) {
	e.GET("/a/:article_id", ShowArticleLanding(service, apRoot))
}

// ShowArticleLanding 共有URLの着地ページ。
// Cache-Controlはここで自己管理する(CacheControlミドルウェアはステータス別の
// 出し分けができないため/a/を対象にしない)。記事は不変なのでCFエッジに1日置く
func ShowArticleLanding(service article.UseCase, apRoot string) echo.HandlerFunc {
	return func(c echo.Context) error {
		articleID, err := entity.StringToID(c.Param("article_id"))
		if err != nil {
			return landingNotFound(c)
		}
		foundArticle, err := service.GetArticle(articleID)
		if err == entity.ErrNotFound {
			return landingNotFound(c)
		}
		if err != nil {
			log.Infof("サービスGetArticleが失敗しました: %v", err)
			return landingNotFound(c)
		}
		articleJson, err := presenter.PickArticle(foundArticle)
		if err != nil {
			log.Infof("PickArticleが失敗しました: %v", err)
			return landingNotFound(c)
		}
		view := landingView{
			ID:           articleJson.ID,
			Title:        articleJson.Title,
			SiteTitle:    articleJson.SiteTitle,
			ImageURL:     articleJson.ImageURL,
			CanonicalURL: apRoot + "/a/" + articleJson.ID,
		}
		c.Response().Header().Set("Cache-Control", "public, max-age=0, s-maxage=86400")
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return landingArticleTmpl.Execute(c.Response(), view)
	}
}

func landingNotFound(c echo.Context) error {
	// 404をCFにキャッシュさせない(後からクロールされた記事が開けるように)
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.HTML(http.StatusNotFound, landingNotFoundHTML)
}
