package entity

import (
	"time"

	"github.com/labstack/gommon/log"
)

type Article struct {
	ID          ID
	Title       ArticleTitle
	URL         string
	ImageURL    ImageURL
	Site        Site
	PublishedAt time.Time
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// NewArticle create a new article
func NewArticle(title string, URL string, imageURL ImageURL, site Site, publishedAt time.Time) (Article, error) {
	article := Article{
		ID:          NewID(),
		Title:       NewArticleTitle(title),
		URL:         URL,
		ImageURL:    imageURL,
		Site:        site,
		PublishedAt: publishedAt,
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := article.Validate(); err != nil {
		return Article{}, ErrInvalidEntity
	}
	return article, nil
}

// Validate validate data
func (a *Article) Validate() error {
	if a.Title.Validate() != nil || a.URL == "" || a.ImageURL.Validate() != nil || a.PublishedAt.IsZero() {
		log.Infof("Title: %v, URL: %v, ImageURL: %v, PublishedAt: %v", a.Title.Validate() != nil, a.URL == "", a.ImageURL.Validate() != nil, a.PublishedAt.IsZero())
		return ErrInvalidEntity
	}
	return nil
}
