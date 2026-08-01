package entity

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
)

// FeedParser abstracts RSS feed fetching/parsing (replaced with a stub in tests)
type FeedParser interface {
	ParseURL(feedURL string) (*gofeed.Feed, error)
}

type Stock struct {
	ArticleSet ArticleSet
	SiteList   []Site
	feedParser FeedParser
}

// NewStock create a new entity Stock
func NewStock(siteList []Site) Stock {
	return NewStockWithParser(siteList, gofeed.NewParser())
}

// NewStockWithParser create a new entity Stock with a custom feed parser
func NewStockWithParser(siteList []Site, parser FeedParser) Stock {
	return Stock{
		ArticleSet: NewArticleSet(),
		SiteList:   siteList,
		feedParser: parser,
	}
}

// StockLatestArticle
// 1サイト・1記事の失敗が全サイトのクロールを止めないよう、失敗はスキップして継続する
func (s Stock) StockLatestArticle() ([]Site, ArticleSet, error) {
	var newSiteList []Site
	newArticleSet := NewArticleSet()
	for _, site := range s.SiteList {
		newSite, articleList, ok := s.stockSite(site)
		if !ok {
			continue
		}
		newSiteList = append(newSiteList, newSite)
		for _, article := range articleList {
			newArticleSet = newArticleSet.Add(article)
		}
	}
	return newSiteList, newArticleSet, nil
}

// stockSite 1サイト分の新着記事を取り込む。取得失敗・日付なし・新着なしは ok=false
func (s Stock) stockSite(site Site) (Site, []Article, bool) {
	feed, err := s.feedParser.ParseURL(site.RSSURL)
	if err != nil {
		log.Warnf("%sのParseURLに失敗しました(このサイトのみスキップ) :%v", site.RSSURL, err)
		return Site{}, nil, false
	}
	lastUpdatedAt, ok := latestPublishedAt(feed)
	if !ok {
		log.Warnf("%sのフィードに日付付き記事がありません(このサイトのみスキップ)", site.RSSURL)
		return Site{}, nil, false
	}
	if site.IsNewerLastUpdatedAt(lastUpdatedAt) {
		return Site{}, nil, false
	}
	var articleList []Article
	for _, item := range feed.Items {
		publishedAt := itemPublishedAt(item)
		if publishedAt == nil {
			continue
		}
		if site.IsNewerLastUpdatedAt(*publishedAt) {
			continue
		}
		if NewArticleTitle(item.Title).ContainsBlacklist() {
			continue
		}
		article, err := NewArticle(item.Title, item.Link, itemImageURL(item), site, *publishedAt)
		if err != nil {
			log.Warnf("%sのNewArticleに失敗しました(この記事のみスキップ) :%v", item.Link, err)
			continue
		}
		articleList = append(articleList, article)
	}
	return site.UpdateLastUpdatedAt(lastUpdatedAt), articleList, true
}

// latestPublishedAt フィード内の最新記事日時。先頭が最新とは限らないため全記事を走査する
func latestPublishedAt(feed *gofeed.Feed) (time.Time, bool) {
	var latest time.Time
	found := false
	for _, item := range feed.Items {
		publishedAt := itemPublishedAt(item)
		if publishedAt == nil {
			continue
		}
		if !found || publishedAt.After(latest) {
			latest = *publishedAt
			found = true
		}
	}
	return latest, found
}

// itemPublishedAt 公開日時。RSS1.0のdc:date/RSS2.0のpubDateはPublishedParsedに入るが、
// 更新日時しか持たないフィードのためUpdatedParsedもフォールバックとして見る
func itemPublishedAt(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	}
	return item.UpdatedParsed
}

// itemImageURL サムネイル抽出。content:encoded(livedoor/FC2)→description(WordPress等)の
// 順に本文中のimgを探し、見つからなければフォールバック画像を採用する
func itemImageURL(item *gofeed.Item) ImageURL {
	for _, content := range []string{item.Content, item.Description} {
		if content == "" {
			continue
		}
		imageURL, err := ContentToImageURL(content)
		if err != nil {
			continue
		}
		if _, err := imageURL.URL(); err != nil {
			continue
		}
		// URL()はimgが実在した場合のみValueを設定する。未設定=img無しなので次の候補へ
		if imageURL.Value != "" {
			return imageURL
		}
	}
	fallback, _ := NewImageURL(RandomImage())
	return fallback
}
