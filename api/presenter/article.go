package presenter

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type ArticleResponce struct {
	Data            []*Article `json:"data"`
	LastPublishedAt time.Time  `json:"lastPublishedAt"`
}

type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"titles"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"image"`
	SiteTitle   string    `json:"sitetitle"`
	SiteID      string    `json:"siteID"`
	PublishedAt time.Time `json:"publishedAt"`
}

type ArticleList []*Article

// pickArticle
func pickArticle(article entity.Article) (Article, error) {
	imageURL, err := article.ImageURL.URL()
	if err != nil {
		return Article{}, err
	}
	articlePresenter := Article{
		ID:          article.ID.String(),
		Title:       article.Title.String(),
		URL:         article.URL,
		ImageURL:    imageURL,
		SiteTitle:   article.Site.Title,
		SiteID:      article.Site.ID.String(),
		PublishedAt: article.PublishedAt,
	}
	return articlePresenter, nil
}

// PickArticleList
func PickArticleList(articleList []entity.Article) (ArticleList, error) {
	var articlePresenterList ArticleList
	for _, article := range articleList {
		var articlePresenter Article
		articlePresenter, err := pickArticle(article)
		if err != nil {
			return ArticleList{}, err
		}
		articlePresenterList = append(articlePresenterList, &articlePresenter)
	}
	return articlePresenterList, nil
}
