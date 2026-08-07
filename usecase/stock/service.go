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
		if err := s.siteService.UpdateSite(site); err != nil {
			// 1サイトの更新失敗で他サイトの取り込みを止めない
			log.Infof("サービスUpdateSiteに失敗しました(スキップして継続): site=%v, err=%v", site.Title, err)
			continue
		}
	}
	for _, article := range newArticleSet.Set {
		if err := s.articleService.CreateArticle(article); err != nil {
			// 1記事のINSERT失敗(カラム長超過など)で同サイクルの残り記事を捨てない
			log.Infof("サービスCreateArticleに失敗しました(スキップして継続): title=%v, err=%v", article.Title.String(), err)
			continue
		}
	}
	return nil
}
