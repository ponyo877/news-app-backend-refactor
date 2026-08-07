package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/user"
)

// MakeUserHandlers
// GET /v1/user は全ユーザーのdevicehash(コメントの本人判定に使う値)を無認証で
// 公開していたため閉鎖した。アプリはPOSTしか使っていない
func MakeUserHandlers(e *echo.Echo, service user.UseCase) {
	e.POST("/v1/user", CreateUser(service)) // name, devicehash, avatarURL
}

// CreateUser
func CreateUser(service user.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		deviceHash := c.FormValue("devicehash")
		name := c.FormValue("name")
		avatar, err := c.FormFile("avatar")

		var avatarImage entity.Image
		if err == nil {
			avatarFile, err := avatar.Open()
			if err != nil {
				log.Infof("アップロードされたファイルが開けません: %v", err)
				return c.JSON(http.StatusBadRequest, nil)
			}
			imageByte, err := ioutil.ReadAll(avatarFile)
			if err != nil {
				log.Infof("アップロードされたファイル(%v)を読み込むことができません: %v", avatar.Filename, err)
				return c.JSON(http.StatusBadRequest, nil)
			}
			avatarImage = entity.Image{
				File: imageByte,
				Name: avatar.Filename,
			}
		} else {
			log.Infof("パラメータavatarが適切に指定されていません: %v", err)
			avatarImage = entity.Image{}
		}

		if _, err := service.CreateUser(name, avatarImage, deviceHash); err != nil {
			log.Infof("サービスCreateUserが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
