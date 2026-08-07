package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// CacheControl は読み取り系GETにCache-Controlを付与し、Cloudflareエッジにのみキャッシュさせる。
// 全リクエストが単一VMのMySQLに直撃する構造を避けるための前提。
//
// max-age=0 + s-maxage の組み合わせが重要:
//   - s-maxage: Cloudflareエッジのキャッシュ時間(オリジン保護)
//   - max-age=0: 端末側(iOSのNSURLCache等)にはキャッシュさせない。
//     max-ageを正の値にすると、アプリの新着リストが端末内キャッシュで数分〜数時間止まって見える
//
// 注意: Cloudflareダッシュボード側の前提設定が2つある
//  1. Cache Rule「パス /v1/* を Eligible for cache(Origin Cache Control尊重)」
//  2. Caching → Configuration → Browser Cache TTL を「既存のヘッダーを尊重する」
//     (既定の4時間のままだと max-age が14400に書き換えられ、端末キャッシュが発生する)
func CacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodGet {
			if maxAge := cacheMaxAge(c.Request().URL.Path); maxAge > 0 {
				c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=0, s-maxage=%d", maxAge))
			}
		}
		return next(c)
	}
}

// cacheMaxAge はパスごとのキャッシュ秒数。0はキャッシュさせない。
func cacheMaxAge(path string) int {
	switch {
	case path == "/v1/article":
		// 新着はクロール周期(2分)より短くして鮮度を保つ
		return 60
	case strings.HasPrefix(path, "/v1/article/view/popular/"):
		return 300
	case path == "/v1/article/search":
		return 300
	case path == "/v1/article/recommend":
		return 300
	case strings.HasPrefix(path, "/v1/article/similar/"):
		return 3600
	case path == "/v1/site":
		return 3600
	}
	return 0
}

// CronGuard は内部バッチ専用エンドポイント(/v1/stock 系)を共有シークレットで保護する。
// 呼び出し側(systemdタイマーのcurl)は -H "X-Cron-Token: $CRON_TOKEN" を付ける。
// CRON_TOKEN未設定時は互換のため従来どおり通す(デプロイ直後にクロールを止めないため)。
// 未設定のままだと外部から誰でもクロール・BERT索引再構築を起動できるので、必ず設定すること。
func CronGuard(token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if token == "" {
				log.Warn("CRON_TOKEN未設定のため /v1/stock 系が無認証で公開されています")
				return next(c)
			}
			if c.Request().Header.Get("X-Cron-Token") != token {
				// 存在自体を隠すため404を返す
				return c.NoContent(http.StatusNotFound)
			}
			return next(c)
		}
	}
}
