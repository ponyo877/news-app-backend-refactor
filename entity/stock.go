package entity

import (
	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
)

type Stock struct {
	ArticleSet ArticleSet
	SiteList   []Site
	feedParser *gofeed.Parser
}

// NewStock create a new entity Stock
func NewStock(siteList []Site) Stock {
	return Stock{
		ArticleSet: NewArticleSet(),
		SiteList:   siteList,
		feedParser: gofeed.NewParser(),
	}
}

// StockLatestArticle
func (s Stock) StockLatestArticle() ([]Site, ArticleSet, error) {
	var newSiteList []Site
	newArticleSet := NewArticleSet()
	for _, site := range s.SiteList {
		feed, err := s.feedParser.ParseURL(site.RSSURL)
		if err != nil {
			log.Infof("%sのParseURLに失敗しました :%v", site.RSSURL, err)
			return nil, NewArticleSet(), err
		}
		lastUpdatedAt := *feed.Items[0].PublishedParsed //*feed.PublishedParsed
		if site.IsNewerLastUpdatedAt(lastUpdatedAt) {
			continue
		}
		newSite := site.UpdateLastUpdatedAt(lastUpdatedAt)
		newSiteList = append(newSiteList, newSite)
		for _, item := range feed.Items {
			publishedAt := *item.PublishedParsed
			if site.IsNewerLastUpdatedAt(publishedAt) {
				continue
			}
			if NewArticleTitle(item.Title).ContainsBlacklist() {
				continue
			}
			imageURL, err := ContentToImageURL(item.Content)
			if err != nil {
				log.Infof("ContentToImageURLに失敗しました :%v", err)
				return nil, NewArticleSet(), err
			}
			newArticle, err := NewArticle(item.Title, item.Link, imageURL, site, publishedAt)
			if err != nil {
				log.Infof("NewArticleに失敗しました :%v", err)
				return nil, NewArticleSet(), err
			}
			newArticleSet = newArticleSet.Add(newArticle)
		}
	}
	return newSiteList, newArticleSet, nil
}
