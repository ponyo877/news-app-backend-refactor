package repository

import (
	"bytes"
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
func (r *ImageNginx) Download(filename string) (entity.Image, error) {
	return entity.Image{}, nil
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
