package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/report"
)

// MakeReportHandlers
func MakeReportHandlers(e *echo.Echo, service report.UseCase) {
	e.POST("/v1/article/report/:article_id", ReportArticle(service)) // url, sitetitle, devicehash, platform, appversion, rulesversion, reason
}

// ReportArticle 記事表示の不具合報告(アプリの⋮メニュー)。同一端末×記事の再報告も200
func ReportArticle(service report.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		articleID, err := entity.StringToID(c.Param("article_id"))
		if err != nil {
			log.Infof("パラメータarticle_idの形式が間違っています: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		// 数値でなければ0(未指定と同じ扱い)。ルールのバージョンは付帯情報なので拒否はしない
		rulesVersion, _ := strconv.Atoi(c.FormValue("rulesversion"))
		err = service.ReportArticle(
			articleID,
			c.FormValue("url"),
			c.FormValue("sitetitle"),
			c.FormValue("devicehash"),
			c.FormValue("platform"),
			c.FormValue("appversion"),
			rulesVersion,
			c.FormValue("reason"),
		)
		if errors.Is(err, entity.ErrInvalidEntity) {
			log.Infof("報告の内容が不正です: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		if err != nil {
			// DB障害等。アプリ側は 5xx を「サーバーが混み合っています」と表示し再試行できる
			log.Warnf("サービスReportArticleが失敗しました: %v", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
