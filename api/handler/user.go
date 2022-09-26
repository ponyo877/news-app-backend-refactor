package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/usecase/user"
)

func MakeUserHandlers(e *echo.Echo, service user.UseCase) {
	e.POST("/v1/user", CreateUser(service)) // name, devicehash, avatarURL
}

// CreateUser
func CreateUser(service user.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Responce{
			Data: []*presenter.Article{},
		}
		return c.JSON(http.StatusOK, responce)
	}
}
