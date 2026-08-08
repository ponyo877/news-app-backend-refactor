package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/usecase/notification"
)

// MakeNotificationHandlers
func MakeNotificationHandlers(e *echo.Echo, service notification.UseCase, cronGuard echo.MiddlewareFunc) {
	e.POST("/v1/user/token", RegisterDeviceToken(service))            // expotoken, devicehash, platform, digest
	e.POST("/v1/notification/digest", SendDigest(service), cronGuard) // ダイジェスト送信(systemdタイマー専用)
}

// RegisterDeviceToken
func RegisterDeviceToken(service notification.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		expoToken := c.FormValue("expotoken")
		deviceHash := c.FormValue("devicehash")
		platform := c.FormValue("platform")
		// 未指定は許諾直後の初回登録なのでON扱い。明示的な "false" のみOFF
		digestEnabled := c.FormValue("digest") != "false"
		if err := service.RegisterToken(expoToken, deviceHash, platform, digestEnabled); err != nil {
			log.Infof("サービスRegisterTokenが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}

// SendDigest
func SendDigest(service notification.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		sentCount, err := service.SendDailyDigest()
		if err != nil {
			log.Infof("サービスSendDailyDigestが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		log.Infof("ダイジェスト通知を%d件送信しました", sentCount)
		return c.JSON(http.StatusOK, map[string]int{"sent": sentCount})
	}
}
