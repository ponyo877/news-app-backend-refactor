package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/comment"
)

// MakeCommentHandlers
func MakeCommentHandlers(e *echo.Echo, service comment.UseCase) {
	e.GET("/v1/comment/:article_id", ListComment(service))    // articleID
	e.POST("/v1/comment/:article_id", CreateComment(service)) // articleID, massage, devicehash
}

// ListComment
func ListComment(service comment.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		articleIDString := c.Param("article_id")
		articleID, err := entity.StringToID(articleIDString)
		if err != nil {
			log.Infof("パラメータarticle_idの形式が間違っています: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		comments, err := service.ListComments(articleID)
		if err != nil {
			log.Infof("サービスListCommentsが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		commentJson, err := presenter.PickCommentList(comments)
		if err != nil {
			log.Infof("PickCommentListが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		responce := presenter.CommentResponce{
			Data: commentJson,
		}
		return c.JSON(http.StatusOK, responce)
	}
}

// CreateComment
func CreateComment(service comment.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		articleIDString := c.Param("article_id")
		articleID, err := entity.StringToID(articleIDString)
		if err != nil {
			log.Infof("パラメータarticle_idの形式が間違っています: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		message := c.FormValue("message")
		deviceHash := c.FormValue("devicehash")

		comment := entity.Comment{
			ID:         entity.NewID(),
			UserName:   "",
			AvatarURL:  "",
			DeviceHash: deviceHash,
			Message:    message,
			Article: entity.Article{
				ID: articleID,
			},
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}
		if _, err := service.CreateComment(comment); err != nil {
			log.Infof("サービスCreateCommentが失敗しました: %v", err)
			c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
