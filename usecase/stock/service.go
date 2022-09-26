package stock

import (
	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
	"github.com/ponyo877/news-app-backend-refactor/usecase/site"
)

// Service Article usecase
type Service struct {
	articleService article.UseCase
	siteService    site.UseCase
	feedParser     *gofeed.Parser
}

// NewService create new service
func NewService(a article.UseCase, s site.UseCase) *Service {
	feedParser := gofeed.NewParser()
	return &Service{
		articleService: a,
		siteService:    s,
		feedParser:     feedParser,
	}
}

// StockLatestArticle stock a laterst article
func (s *Service) StockLatestArticle() error {
	siteList, err := s.siteService.ListSite()
	if err != nil {
		return err
	}
	stock := entity.NewStock(siteList)
	newSiteList, newArticleSet, err := stock.StockLatestArticle()
	if err != nil {
		log.Infof("エンティティメソッドStockLatestArticleに失敗しました: %v", err)
		return err
	}
	for _, site := range newSiteList {
		if err := s.siteService.UpdateSite(&site); err != nil {
			log.Infof("サービスUpdateSiteに失敗しました: %v", err)
			return err
		}
	}
	for _, article := range newArticleSet.Set {
		if _, _, err := s.articleService.CreateArticle(article); err != nil {
			log.Infof("サービスCreateArticleに失敗しました: %v", err)
			return err
		}
	}
	return nil
}
