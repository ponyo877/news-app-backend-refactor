package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/api/presenter"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/user"
)

// MakeUserHandlers
func MakeUserHandlers(e *echo.Echo, service user.UseCase) {
	e.GET("/v1/user", ListUsers(service))
	e.POST("/v1/user", CreateUser(service)) // name, devicehash, avatarURL
}

// ListUsers
func ListUsers(service user.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := service.ListUsers()
		if err == entity.ErrNotFound {
			return c.JSON(http.StatusOK, presenter.UserResponce{
				Data: []*presenter.User{},
			})
		}
		if err != nil {
			log.Infof("サービスListUserが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		userJson, err := presenter.PickUserList(users)
		if err != nil {
			log.Infof("PickUserListが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		responce := presenter.UserResponce{
			Data: userJson,
		}
		return c.JSON(http.StatusOK, responce)
	}
}

// CreateUser
func CreateUser(service user.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		deviceHash := c.FormValue("devicehash")
		avatar, err := c.FormFile("avatar")
		if err != nil {
			log.Infof("パラメータavatarの形式が間違っています: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		avatarFile, err := avatar.Open()
		if err != nil {
			log.Infof("アップロードされたファイルが開けません: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		imageByte, err := ioutil.ReadAll(avatarFile)
		if err != nil {
			log.Infof("パラメータavatarの形式が間違っています: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		avatarImage := entity.Image{
			File: imageByte,
			Name: avatar.Filename,
		}
		if _, err := service.CreateUser(name, avatarImage, deviceHash); err != nil {
			log.Infof("サービスCreateUserが失敗しました: %v", err)
			return c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
