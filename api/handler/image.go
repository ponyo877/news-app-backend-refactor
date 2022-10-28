package handler

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/fileio"
)

func MakeImageHandlers(e *echo.Echo, service fileio.UseCase) {
	e.GET("/v1/static/:filename", FetchImage(service))
	e.POST("/v1/static", SaveImage(service))
}

// FetchImage
func FetchImage(service fileio.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		filename := c.Param("filename")
		image, err := service.FetchImage(filename)
		if err != nil {
			log.Infof("サービスFetchImageが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		io.Copy(c.Response().Writer, bytes.NewReader(image.File))
		return c.NoContent(http.StatusOK)
	}
}

// SaveImage
func SaveImage(service fileio.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			log.Infof("パラメータfileの形式が間違っています: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		imageFile, err := file.Open()
		if err != nil {
			log.Infof("アップロードされたファイルが開けません: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		imageByte, err := ioutil.ReadAll(imageFile)
		if err != nil {
			log.Infof("パラメータfileの形式が間違っています: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		image := entity.Image{
			File: imageByte,
			Name: file.Filename,
		}
		if _, err := service.SaveImage(image); err != nil {
			log.Infof("サービスSaveImageが失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
