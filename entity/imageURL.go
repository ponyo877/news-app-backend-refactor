package entity

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ImageURL struct {
	Content string
	Value   string
}

// NewComment create a new article
func NewImageURL(imageURL string) (ImageURL, error) {
	newImageURL := ImageURL{
		Value: imageURL,
	}
	if err := newImageURL.Validate(); err != nil {
		return ImageURL{}, ErrInvalidEntity
	}
	return newImageURL, nil
}

// ContentToImangeURL
func ContentToImangeURL(content string) (ImageURL, error) {
	newImageURL := ImageURL{
		Content: content,
	}
	if err := newImageURL.Validate(); err != nil {
		return ImageURL{}, ErrInvalidEntity
	}
	return newImageURL, nil
}

// Validate validate data
func (i *ImageURL) Validate() error {
	if i.Value == "" && i.Content == "" {
		return ErrInvalidEntity
	}
	return nil
}

// URL
func (i *ImageURL) URL() (string, error) {
	if i.Value != "" {
		return i.Value, nil
	}
	reader := strings.NewReader(i.Content)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}
	imageUrl, exist := doc.Find("img").Attr("src")
	if !exist {
		return "", ErrNotFound
	}
	i.Value = imageUrl
	return imageUrl, nil
}
