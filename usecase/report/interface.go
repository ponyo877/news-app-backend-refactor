package report

import "github.com/ponyo877/news-app-backend-refactor/entity"

// Repository interface
type Repository interface {
	Save(e entity.ArticleReport) error
}

// UseCase interface
type UseCase interface {
	ReportArticle(articleID entity.ID, url, siteTitle, deviceHash, platform, appVersion string, rulesVersion int, reason string) error
}
