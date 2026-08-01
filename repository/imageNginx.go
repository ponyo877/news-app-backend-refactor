package repository

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// ImageNginx webdav repository
type ImageNginx struct {
	endpoint string
	dir      string
}

// NewImageNginx create new repository
func NewImageNginx(endpoint string) *ImageNginx {
	return &ImageNginx{
		endpoint: endpoint,
		dir:      "static",
	}
}

// Download
// 従来は未実装スタブで、GET /v1/static/:filename が常に0バイトを返していた。
// nginx(WebDAV)の /static から実体を取得して返す
func (r *ImageNginx) Download(filename string) (entity.Image, error) {
	webdavURL, err := url.Parse(r.endpoint)
	if err != nil {
		return entity.Image{}, err
	}
	webdavURL.Path = filepath.Join(webdavURL.Path, r.dir, filename)
	res, err := http.Get(webdavURL.String())
	if err != nil {
		return entity.Image{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return entity.Image{}, fmt.Errorf("静的ファイルの取得に失敗しました: HTTP %d", res.StatusCode)
	}
	file, err := io.ReadAll(res.Body)
	if err != nil {
		return entity.Image{}, err
	}
	return entity.Image{Name: filename, File: file}, nil
}

// Upload
func (r *ImageNginx) Upload(e entity.Image) (string, error) {
	reader := bytes.NewReader(e.File)
	webdavURL, err := url.Parse(r.endpoint)
	if err != nil {
		return "", err
	}
	webdavURL.Path = filepath.Join(webdavURL.Path, r.dir, e.FileName())
	req, err := http.NewRequest("PUT", webdavURL.String(), reader)
	if err != nil {
		return "", err
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		return "", err
	}
	return filepath.Join(r.dir, e.FileName()), nil
}
