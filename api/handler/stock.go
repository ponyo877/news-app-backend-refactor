package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/usecase/stock"
)

// MakeStockHandlers
func MakeStockHandlers(e *echo.Echo, service stock.UseCase) {
	e.GET("/v1/stock", StockLatestArticle(service)) // 全サイトの情報を得る
}

// StockLatestArticle
func StockLatestArticle(service stock.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := service.StockLatestArticle(); err != nil {
			log.Infof("サービスStockLatestArticleが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		log.Info("サービスStockLatestArticleが成功しました")
		return c.JSON(http.StatusOK, nil)
	}
}
