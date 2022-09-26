package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/usecase/site"
)

func MakeSiteHandlers(e *echo.Echo, service site.UseCase) {
	e.GET("/v1/site", ListSite(service)) // 全サイトの情報を得る
}

// ListSite
func ListSite(service site.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Responce{
			Data: []*presenter.Article{},
		}
		return c.JSON(http.StatusOK, responce)
	}
}
