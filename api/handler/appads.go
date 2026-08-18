package handler

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

// app-ads.txt (IAB Tech Lab / Authorized Sellers for Apps)。
// ストア掲載の「デベロッパーのウェブサイト」に指定したドメインのルートに置く必要があり、
// AdMobのクローラーがそこを見て正規の販売者かを判定する。
// 未設置でも広告配信は止まらないが、なりすまし在庫による収益の目減りを防ぐ。
//
// テンプレート類(landing_article.html)と同じくembedでバイナリに同梱し、
// CGO_ENABLED=0のクロスコンパイルに影響させない
//
//go:embed app-ads.txt
var appAdsTxt string

// MakeAppAdsHandlers
// HEADも受ける(クローラーがGETの前にHEADを打つことがある)
func MakeAppAdsHandlers(e *echo.Echo) {
	h := showAppAdsTxt()
	e.GET("/app-ads.txt", h)
	e.HEAD("/app-ads.txt", h)
}

// ShowAppAdsTxt text/plainで返す(仕様上、他のContent-Typeは無視される)
func showAppAdsTxt() echo.HandlerFunc {
	return func(c echo.Context) error {
		// 内容はほぼ不変だが、変更を1日以上引きずらせないよう短めに置く
		c.Response().Header().Set("Cache-Control", "public, max-age=0, s-maxage=3600")
		return c.String(http.StatusOK, appAdsTxt)
	}
}
