package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/notification"
)

// MakeNotificationHandlers
func MakeNotificationHandlers(e *echo.Echo, service notification.UseCase, cronGuard echo.MiddlewareFunc) {
	e.POST("/v1/user/token", RegisterDeviceToken(service))              // expotoken, devicehash, platform, digest, matsuri
	e.POST("/v1/notification/digest", SendDigest(service), cronGuard)   // ダイジェスト送信(systemdタイマー専用)
	e.POST("/v1/notification/matsuri", SendMatsuri(service), cronGuard) // 祭り速報(Cloudflare Worker専用)
}

// RegisterDeviceToken
func RegisterDeviceToken(service notification.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		expoToken := c.FormValue("expotoken")
		deviceHash := c.FormValue("devicehash")
		platform := c.FormValue("platform")
		// 未指定は許諾直後の初回登録なのでON扱い。明示的な "false" のみOFF
		digestEnabled := c.FormValue("digest") != "false"
		matsuriEnabled := c.FormValue("matsuri") != "false"
		if err := service.RegisterToken(expoToken, deviceHash, platform, digestEnabled, matsuriEnabled); err != nil {
			log.Infof("サービスRegisterTokenが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}

// matsuriRequest 祭り検知元(Worker)が送る記事メタ+サイト数。キー体系は一覧APIと同一
type matsuriRequest struct {
	ID          string `json:"id"`
	Titles      string `json:"titles"`
	URL         string `json:"url"`
	Image       string `json:"image"`
	SiteID      string `json:"siteID"`
	SiteTitle   string `json:"sitetitle"`
	PublishedAt string `json:"publishedAt"`
	SiteCount   int    `json:"siteCount"`
}

// SendMatsuri
func SendMatsuri(service notification.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req matsuriRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		articleID, err := entity.StringToID(req.ID)
		if err != nil || req.Titles == "" || req.SiteCount < 2 {
			return c.JSON(http.StatusBadRequest, nil)
		}
		siteID, err := entity.StringToID(req.SiteID)
		if err != nil {
			siteID = entity.NewID()
		}
		publishedAt, err := time.Parse(time.RFC3339, req.PublishedAt)
		if err != nil {
			publishedAt = time.Now()
		}
		matsuriArticle := entity.Article{
			ID:          articleID,
			Title:       entity.NewArticleTitle(req.Titles),
			URL:         req.URL,
			Site:        entity.Site{ID: siteID, Title: req.SiteTitle},
			PublishedAt: publishedAt,
		}
		sentCount, err := service.SendMatsuri(matsuriArticle, req.Image, req.SiteCount)
		if err != nil {
			log.Infof("サービスSendMatsuriが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		log.Infof("祭り速報を%d件送信しました: %s", sentCount, req.Titles)
		return c.JSON(http.StatusOK, map[string]int{"sent": sentCount})
	}
}

// SendDigest
func SendDigest(service notification.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		sentCount, err := service.SendDailyDigest()
		if err != nil {
			log.Infof("サービスSendDailyDigestが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		log.Infof("ダイジェスト通知を%d件送信しました", sentCount)
		return c.JSON(http.StatusOK, map[string]int{"sent": sentCount})
	}
}
