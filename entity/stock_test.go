package entity

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

type stubParser struct {
	feeds map[string]*gofeed.Feed
	errs  map[string]error
}

func (s *stubParser) ParseURL(feedURL string) (*gofeed.Feed, error) {
	if err, ok := s.errs[feedURL]; ok {
		return nil, err
	}
	feed, ok := s.feeds[feedURL]
	if !ok {
		return nil, errors.New("unexpected url: " + feedURL)
	}
	return feed, nil
}

func testSite(t *testing.T, title, rssURL string, lastUpdatedAt time.Time) Site {
	t.Helper()
	site, err := NewSite(title, rssURL, "https://example.com/icon.png")
	if err != nil {
		t.Fatalf("NewSite failed: %v", err)
	}
	return site.UpdateLastUpdatedAt(lastUpdatedAt)
}

func feedItem(title, link string, publishedAt *time.Time, content string) *gofeed.Item {
	return &gofeed.Item{Title: title, Link: link, PublishedParsed: publishedAt, Content: content}
}

func timePtr(t time.Time) *time.Time { return &t }

var (
	t1 = time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	t2 = time.Date(2026, 8, 1, 11, 0, 0, 0, time.UTC)
	t3 = time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
)

// 1サイトのRSS取得失敗が他サイトの取り込みを止めないこと
func TestStockLatestArticleContinuesWhenOneSiteFails(t *testing.T) {
	parser := &stubParser{
		feeds: map[string]*gofeed.Feed{
			"https://a.example.com/rss": {Items: []*gofeed.Item{feedItem("記事A", "https://a.example.com/1", timePtr(t1), "")}},
			"https://c.example.com/rss": {Items: []*gofeed.Item{feedItem("記事C", "https://c.example.com/1", timePtr(t2), "")}},
		},
		errs: map[string]error{"https://b.example.com/rss": errors.New("connection refused")},
	}
	stock := NewStockWithParser([]Site{
		testSite(t, "サイトA", "https://a.example.com/rss", time.Time{}),
		testSite(t, "サイトB", "https://b.example.com/rss", time.Time{}),
		testSite(t, "サイトC", "https://c.example.com/rss", time.Time{}),
	}, parser)

	newSiteList, articleSet, err := stock.StockLatestArticle()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(newSiteList) != 2 {
		t.Fatalf("newSiteList = %d sites, want 2 (failed site must be skipped)", len(newSiteList))
	}
	if len(articleSet.Set) != 2 {
		t.Fatalf("articles = %d, want 2", len(articleSet.Set))
	}
}

// 空フィード・日付なしフィードでpanicせずスキップされること
func TestStockLatestArticleSkipsEmptyOrDatelessFeeds(t *testing.T) {
	parser := &stubParser{
		feeds: map[string]*gofeed.Feed{
			"https://empty.example.com/rss":    {Items: []*gofeed.Item{}},
			"https://dateless.example.com/rss": {Items: []*gofeed.Item{feedItem("日付なし", "https://dateless.example.com/1", nil, "")}},
			"https://ok.example.com/rss":       {Items: []*gofeed.Item{feedItem("記事", "https://ok.example.com/1", timePtr(t1), "")}},
		},
	}
	stock := NewStockWithParser([]Site{
		testSite(t, "空サイト", "https://empty.example.com/rss", time.Time{}),
		testSite(t, "日付なしサイト", "https://dateless.example.com/rss", time.Time{}),
		testSite(t, "正常サイト", "https://ok.example.com/rss", time.Time{}),
	}, parser)

	newSiteList, articleSet, err := stock.StockLatestArticle()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(newSiteList) != 1 || len(articleSet.Set) != 1 {
		t.Fatalf("sites=%d articles=%d, want 1/1", len(newSiteList), len(articleSet.Set))
	}
}

// PublishedParsedが無い記事はUpdatedParsedで代替されること
func TestStockLatestArticleFallsBackToUpdatedParsed(t *testing.T) {
	item := &gofeed.Item{Title: "更新日時のみ", Link: "https://u.example.com/1", UpdatedParsed: timePtr(t2)}
	parser := &stubParser{feeds: map[string]*gofeed.Feed{
		"https://u.example.com/rss": {Items: []*gofeed.Item{item}},
	}}
	stock := NewStockWithParser([]Site{testSite(t, "更新日時サイト", "https://u.example.com/rss", time.Time{})}, parser)

	newSiteList, articleSet, _ := stock.StockLatestArticle()
	if len(articleSet.Set) != 1 {
		t.Fatalf("articles = %d, want 1", len(articleSet.Set))
	}
	if !newSiteList[0].LastUpdatedAt.Equal(t2) {
		t.Fatalf("LastUpdatedAt = %v, want %v", newSiteList[0].LastUpdatedAt, t2)
	}
}

// watermark(last_updated_at)より新しい記事のみ取り込まれ、watermarkが最新記事時刻に進むこと。
// フィードの先頭が最新でなくても最新時刻を検出できること
func TestStockLatestArticleRespectsWatermarkAndUnsortedFeed(t *testing.T) {
	parser := &stubParser{feeds: map[string]*gofeed.Feed{
		"https://w.example.com/rss": {Items: []*gofeed.Item{
			feedItem("古い記事", "https://w.example.com/1", timePtr(t1), ""),
			feedItem("最新記事", "https://w.example.com/3", timePtr(t3), ""),
			feedItem("途中の記事", "https://w.example.com/2", timePtr(t2), ""),
		}},
	}}
	stock := NewStockWithParser([]Site{testSite(t, "サイトW", "https://w.example.com/rss", t1)}, parser)

	newSiteList, articleSet, _ := stock.StockLatestArticle()
	if len(articleSet.Set) != 2 {
		t.Fatalf("articles = %d, want 2 (t1以前は除外)", len(articleSet.Set))
	}
	if !newSiteList[0].LastUpdatedAt.Equal(t3) {
		t.Fatalf("LastUpdatedAt = %v, want %v", newSiteList[0].LastUpdatedAt, t3)
	}
}

// 新着が無いサイトはスキップされること(境界: 同時刻は新着扱いしない)
func TestStockLatestArticleSkipsSiteWithoutNewArticles(t *testing.T) {
	parser := &stubParser{feeds: map[string]*gofeed.Feed{
		"https://s.example.com/rss": {Items: []*gofeed.Item{feedItem("既読記事", "https://s.example.com/1", timePtr(t2), "")}},
	}}
	stock := NewStockWithParser([]Site{testSite(t, "サイトS", "https://s.example.com/rss", t2)}, parser)

	newSiteList, articleSet, _ := stock.StockLatestArticle()
	if len(newSiteList) != 0 || len(articleSet.Set) != 0 {
		t.Fatalf("sites=%d articles=%d, want 0/0", len(newSiteList), len(articleSet.Set))
	}
}

// サムネイル: content:encodedのimg → descriptionのimg → フォールバック画像 の優先順で解決すること
func TestStockLatestArticleResolvesThumbnail(t *testing.T) {
	items := []*gofeed.Item{
		{Title: "contentに画像", Link: "https://i.example.com/1", PublishedParsed: timePtr(t1),
			Content: `<img src="https://img.example.com/content.jpg" />`},
		{Title: "descriptionに画像", Link: "https://i.example.com/2", PublishedParsed: timePtr(t2),
			Content: `<p>画像なし</p>`, Description: `<img src="https://img.example.com/desc.jpg" />`},
		{Title: "画像なし", Link: "https://i.example.com/3", PublishedParsed: timePtr(t3)},
	}
	parser := &stubParser{feeds: map[string]*gofeed.Feed{
		"https://i.example.com/rss": {Items: items},
	}}
	stock := NewStockWithParser([]Site{testSite(t, "画像サイト", "https://i.example.com/rss", time.Time{})}, parser)

	_, articleSet, _ := stock.StockLatestArticle()
	if len(articleSet.Set) != 3 {
		t.Fatalf("articles = %d, want 3", len(articleSet.Set))
	}
	got := map[string]string{}
	for _, article := range articleSet.Set {
		url, err := article.ImageURL.URL()
		if err != nil {
			t.Fatalf("ImageURL.URL() failed: %v", err)
		}
		got[article.URL] = url
	}
	if got["https://i.example.com/1"] != "https://img.example.com/content.jpg" {
		t.Errorf("content:encoded優先のはず: %s", got["https://i.example.com/1"])
	}
	if got["https://i.example.com/2"] != "https://img.example.com/desc.jpg" {
		t.Errorf("descriptionフォールバックのはず: %s", got["https://i.example.com/2"])
	}
	if !strings.HasPrefix(got["https://i.example.com/3"], "https://matome.folks-chat.com/static/myimage_") {
		t.Errorf("フォールバック画像のはず: %s", got["https://i.example.com/3"])
	}
}

// 同一サイト×同一タイトルの記事が重複排除されること
func TestArticleSetDeduplicatesBySiteAndTitle(t *testing.T) {
	site := testSite(t, "重複サイト", "https://d.example.com/rss", time.Time{})
	articleSet := NewArticleSet()
	for i := 0; i < 2; i++ {
		article, err := NewArticle("同じタイトル", "https://d.example.com/1", mustImage(t), site, t1)
		if err != nil {
			t.Fatalf("NewArticle failed: %v", err)
		}
		articleSet = articleSet.Add(article)
	}
	if len(articleSet.Set) != 1 {
		t.Fatalf("articles = %d, want 1 (dedup)", len(articleSet.Set))
	}
}

func mustImage(t *testing.T) ImageURL {
	t.Helper()
	imageURL, err := NewImageURL("https://img.example.com/a.jpg")
	if err != nil {
		t.Fatalf("NewImageURL failed: %v", err)
	}
	return imageURL
}

// 実在フィード3形式(livedoor RDF / WordPress RSS2.0 / FC2)のフィクスチャが
// gofeedで日付・本文とも取得でき、取り込みまで通ること
func TestStockLatestArticleParsesRealWorldFeedFormats(t *testing.T) {
	cases := []struct {
		name        string
		fixture     string
		wantItems   int
		wantContent bool // content:encodedが取れる形式か
	}{
		{"livedoor RDF (dc:date)", "testdata/feed_livedoor_rdf.xml", 2, true},
		{"WordPress RSS2.0 (pubDate, content:encodedなし)", "testdata/feed_wordpress_rss2.xml", 2, false},
		{"FC2 RDF (dc:date)", "testdata/feed_fc2_rdf.xml", 2, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, err := os.ReadFile(c.fixture)
			if err != nil {
				t.Fatalf("fixture読み込み失敗: %v", err)
			}
			feed, err := gofeed.NewParser().ParseString(string(raw))
			if err != nil {
				t.Fatalf("ParseString失敗: %v", err)
			}
			if len(feed.Items) != c.wantItems {
				t.Fatalf("items = %d, want %d", len(feed.Items), c.wantItems)
			}
			for _, item := range feed.Items {
				if itemPublishedAt(item) == nil {
					t.Errorf("%s: 日付が取得できない", item.Link)
				}
				if c.wantContent && item.Content == "" {
					t.Errorf("%s: content:encodedが取得できない", item.Link)
				}
			}

			parser := &stubParser{feeds: map[string]*gofeed.Feed{"https://fixture.example.com/rss": feed}}
			stock := NewStockWithParser([]Site{testSite(t, "フィクスチャ", "https://fixture.example.com/rss", time.Time{})}, parser)
			newSiteList, articleSet, err := stock.StockLatestArticle()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(newSiteList) != 1 || len(articleSet.Set) != c.wantItems {
				t.Fatalf("sites=%d articles=%d, want 1/%d", len(newSiteList), len(articleSet.Set), c.wantItems)
			}
		})
	}
}
