package repository

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/studio-b12/gowebdav"
)

// ImageWebDAV webdav repository
type ImageWebDAV struct {
	wd  *gowebdav.Client
	dir string
}

// NewImageWebDAV create new repository
func NewImageWebDAV(wd *gowebdav.Client) *ImageWebDAV {
	return &ImageWebDAV{
		wd:  wd,
		dir: "/v1/static/",
	}
}

func (r *ImageWebDAV) Download(filename string) (entity.Image, error) {
	filePath := filepath.Join(r.dir, filename)
	reader, err := r.wd.ReadStream(filePath)
	if err != nil {
		return entity.Image{}, err
	}
	buffer := new(bytes.Buffer)
	io.Copy(buffer, reader)
	return entity.Image{
		File: buffer.Bytes(),
		Name: filename,
	}, nil
}

func (r *ImageWebDAV) Upload(e entity.Image) (string, error) {
	file := bytes.NewReader(e.File)
	filePath := filepath.Join(r.dir, e.FileName())
	if err := r.wd.WriteStream(filePath, file, 0644); err != nil {
		return "", err
	}
	return filePath, nil
}
