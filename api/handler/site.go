package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/site"
)

func MakeSiteHandlers(e *echo.Echo, service site.UseCase) {
	e.GET("/v1/site", ListSite(service)) // 全サイトの情報を得る
}

// ListSite
func ListSite(service site.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		sites, err := service.ListSite()
		if err == entity.ErrNotFound {
			return c.JSON(http.StatusOK, presenter.SiteResponce{
				Data: []*presenter.Site{},
			})
		}
		if err != nil {
			log.Infof("サービスListSiteが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		siteJson, err := presenter.PickSiteList(sites)
		if err != nil {
			log.Infof("PickArticleListが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		responce := presenter.SiteResponce{
			Data: siteJson,
		}
		return c.JSON(http.StatusOK, responce)
	}
}
