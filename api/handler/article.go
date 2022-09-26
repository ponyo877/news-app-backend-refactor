package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
)

// MakeBookHandlers make url handlers
func MakeArticleHandlers(e *echo.Echo, service article.UseCase) {
	e.GET("/v1/article", ListArticles(service))                             // lastpublished, skipIDs
	e.GET("/v1/article/view/popular/:period", ListPopularArticles(service)) // kind
	e.GET("/v1/article/search", ListSearchedArticles(service))              // words
	e.POST("/v1/article/view/:article_id", IncrementViewCount(service))     // post_id
	e.GET("/v1/article/recommend", ListRecommendArticle(service))           // ids

	// "/mongo/get?lastpublished=" + lastpublished + "&skipIDs=" + _skipIDs;
	// "/mongo/ranking/" + type
	// "/elastic/get?words=" + searchwords
	// 記事の投稿はバッチ
	// "/redis/put/"
	// "/personal?ids=" + ids
	// "/comment/get?articleID=" + articleID
	// "/comment/put"
	// "/site/get"
	// "/user/put"

	// "/eula/"
	// "/privacy_policy/"
	// "/recom/" + postID
}

// ListArticles
func ListArticles(service article.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		baseCreatedAtString := c.QueryParam("lastpublished")
		baseCreatedAt, err := time.Parse(time.RFC3339, baseCreatedAtString)
		if err != nil {
			log.Infof("パラメータlastpublishedの形式が間違っています: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		invisibleSiteIDSetString := c.QueryParam("skipIDs")
		invisibleIDSet, err := entity.StringToIDSet(invisibleSiteIDSetString)
		if err != nil {
			log.Infof("パラメータskipIDsの形式が間違っています: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		articles, err := service.ListArticles(baseCreatedAt, invisibleIDSet)
		if err == entity.ErrNotFound {
			return c.JSON(http.StatusOK, presenter.Responce{
				Data: []*presenter.Article{},
			})
		}
		if err != nil {
			log.Infof("サービスListArticlesが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		articleForJson, err := presenter.PickArticleList(articles)
		if err != nil {
			log.Infof("PickArticleListが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		responce := presenter.Responce{
			Data: articleForJson,
		}
		return c.JSON(http.StatusOK, responce)
	}
}

// ListPopularArticles
func ListPopularArticles(service article.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		period := c.Param("period")
		articles, err := service.ListPopularArticles(period)
		if err == entity.ErrNotFound {
			return c.JSON(http.StatusOK, presenter.Responce{
				Data: []*presenter.Article{},
			})
		}
		if err != nil {
			log.Infof("サービスListPopularArticlesが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		articleForJson, err := presenter.PickArticleList(articles)
		if err != nil {
			log.Infof("PickArticleListが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		responce := presenter.Responce{
			Data: articleForJson,
		}
		return c.JSON(http.StatusOK, responce)
	}
}

// ListSearchedArticles
func ListSearchedArticles(service article.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Responce{
			Data: []*presenter.Article{},
		}
		return c.JSON(http.StatusOK, responce)
	}
}

// IncrementViewCount
func IncrementViewCount(service article.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		articleIDString := c.Param("article_id")
		articleID, err := entity.StringToID(articleIDString)
		if err != nil {
			log.Infof("StringToIDが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		if err := service.IncrementViewCount(articleID); err != nil {
			log.Infof("サービスIncrementViewCountが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		log.Info("サービスIncrementViewCountが成功しました")
		return c.JSON(http.StatusOK, nil)
	}
}

// ListRecommendArticle
func ListRecommendArticle(service article.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Responce{
			Data: []*presenter.Article{},
		}
		return c.JSON(http.StatusOK, responce)
	}
}
