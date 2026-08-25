package report

import "github.com/ponyo877/news-app-backend-refactor/entity"

// Service Service struct
type Service struct {
	repository Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

// ReportArticle 記事表示の不具合報告を保存する(同一端末×記事はupsert)
func (s *Service) ReportArticle(articleID entity.ID, url, siteTitle, deviceHash, platform, appVersion string, rulesVersion int, reason string) error {
	report, err := entity.NewArticleReport(articleID, url, siteTitle, deviceHash, platform, appVersion, rulesVersion, reason)
	if err != nil {
		return err
	}
	return s.repository.Save(report)
}
