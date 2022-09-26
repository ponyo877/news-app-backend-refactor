package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/usecase/comment"
)

func MakeCommentHandlers(e *echo.Echo, service comment.UseCase) {
	e.GET("/v1/comment", ListComment(service))    // 記事ID
	e.POST("/v1/comment", CreateComment(service)) // articleID, massage, devicehash
}

// ListComment
func ListComment(service comment.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Responce{
			Data: []*presenter.Article{},
		}
		return c.JSON(http.StatusOK, responce)
	}
}

// CreateComment
func CreateComment(service comment.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Responce{
			Data: []*presenter.Article{},
		}
		return c.JSON(http.StatusOK, responce)
	}
}
