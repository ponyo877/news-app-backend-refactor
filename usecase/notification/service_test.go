package notification

import (
	"errors"
	"testing"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type fakeRepository struct {
	saved   []entity.DeviceToken
	tokens  []entity.DeviceToken
	deleted []string
}

func (f *fakeRepository) Save(e entity.DeviceToken) error {
	f.saved = append(f.saved, e)
	return nil
}

func (f *fakeRepository) ListDigestEnabled() ([]entity.DeviceToken, error) {
	return f.tokens, nil
}

func (f *fakeRepository) ListMatsuriEnabled() ([]entity.DeviceToken, error) {
	return f.tokens, nil
}

func (f *fakeRepository) Delete(expoToken string) error {
	f.deleted = append(f.deleted, expoToken)
	return nil
}

type fakeDigestLog struct {
	recentArticleIDs []string
	savedArticleIDs  []string
	savedSentCounts  []int
}

func (f *fakeDigestLog) Save(articleID entity.ID, sentCount int) error {
	f.savedArticleIDs = append(f.savedArticleIDs, articleID.String())
	f.savedSentCounts = append(f.savedSentCounts, sentCount)
	return nil
}

func (f *fakeDigestLog) ListRecentArticleIDs(since time.Time) ([]string, error) {
	return f.recentArticleIDs, nil
}

type fakePusher struct {
	pushed        []entity.PushMessage
	invalidTokens []string
}

func (f *fakePusher) Push(messages []entity.PushMessage) ([]string, error) {
	f.pushed = append(f.pushed, messages...)
	return f.invalidTokens, nil
}

type fakeArticleService struct {
	articlesByPeriod map[string][]entity.Article
}

func (f *fakeArticleService) ListPopularArticles(period string) ([]entity.Article, error) {
	articles, ok := f.articlesByPeriod[period]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return articles, nil
}

func (f *fakeArticleService) CreateArticle(entity.Article) error { return nil }
func (f *fakeArticleService) GetArticle(entity.ID) (entity.Article, error) {
	return entity.Article{}, nil
}
func (f *fakeArticleService) SearchArticles(entity.Keyword) ([]entity.Article, error) {
	return nil, nil
}
func (f *fakeArticleService) ListArticles(time.Time, entity.IDSet) ([]entity.Article, error) {
	return nil, nil
}
func (f *fakeArticleService) IncrementViewCount(entity.ID) error { return nil }
func (f *fakeArticleService) ListSimilarArticles(entity.ID) ([]entity.Article, error) {
	return nil, nil
}
func (f *fakeArticleService) UpdateMLIndex() error { return nil }

func newTestArticle(title string) entity.Article {
	return entity.Article{
		ID:    entity.NewID(),
		Title: entity.NewArticleTitle(title),
	}
}

func newTestToken(token string) entity.DeviceToken {
	return entity.DeviceToken{ExpoToken: token, DeviceHash: "hash", Platform: "ios", DigestEnabled: true}
}

func TestRegisterToken(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository, &fakeDigestLog{}, &fakePusher{}, &fakeArticleService{})

	if err := service.RegisterToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]", "devicehash", "ios", true, true); err != nil {
		t.Fatalf("正常なトークン登録が失敗しました: %v", err)
	}
	if len(repository.saved) != 1 {
		t.Fatalf("保存件数が1ではありません: %d", len(repository.saved))
	}

	if err := service.RegisterToken("bogus-token", "devicehash", "ios", true, true); err == nil {
		t.Fatal("不正な形式のトークンが登録できてしまいました")
	}
	if err := service.RegisterToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]", "devicehash", "web", true, true); err == nil {
		t.Fatal("不正なプラットフォームが登録できてしまいました")
	}
}

func TestSendDailyDigest(t *testing.T) {
	repository := &fakeRepository{tokens: []entity.DeviceToken{
		newTestToken("ExponentPushToken[aaa]"),
		newTestToken("ExponentPushToken[bbb]"),
	}}
	pusher := &fakePusher{invalidTokens: []string{"ExponentPushToken[bbb]"}}
	articleService := &fakeArticleService{articlesByPeriod: map[string][]entity.Article{
		"daily": {newTestArticle("今日の1位記事")},
	}}
	service := NewService(repository, &fakeDigestLog{}, pusher, articleService)

	sentCount, err := service.SendDailyDigest()
	if err != nil {
		t.Fatalf("ダイジェスト送信が失敗しました: %v", err)
	}
	if sentCount != 1 {
		t.Fatalf("送信成功数が1(2件送信-1件失効)ではありません: %d", sentCount)
	}
	if len(pusher.pushed) != 2 {
		t.Fatalf("送信メッセージ数が2ではありません: %d", len(pusher.pushed))
	}
	if pusher.pushed[0].Body != "今日の1位記事" {
		t.Fatalf("通知本文が記事タイトルではありません: %s", pusher.pushed[0].Body)
	}
	if pusher.pushed[0].Data["id"] == "" {
		t.Fatal("通知データに記事IDがありません")
	}
	if pusher.pushed[0].Data["type"] != "digest" {
		t.Fatal("通知データにtype=digestがありません")
	}
	if len(repository.deleted) != 1 || repository.deleted[0] != "ExponentPushToken[bbb]" {
		t.Fatalf("失効トークンが削除されていません: %v", repository.deleted)
	}
}

func TestSendDailyDigestFallback(t *testing.T) {
	repository := &fakeRepository{tokens: []entity.DeviceToken{newTestToken("ExponentPushToken[aaa]")}}
	pusher := &fakePusher{}
	// dailyが空(朝7時など)→ weeklyへフォールバック
	articleService := &fakeArticleService{articlesByPeriod: map[string][]entity.Article{
		"weekly": {newTestArticle("今週の1位記事")},
	}}
	service := NewService(repository, &fakeDigestLog{}, pusher, articleService)

	sentCount, err := service.SendDailyDigest()
	if err != nil {
		t.Fatalf("フォールバック付きダイジェスト送信が失敗しました: %v", err)
	}
	if sentCount != 1 {
		t.Fatalf("送信成功数が1ではありません: %d", sentCount)
	}
	if pusher.pushed[0].Body != "今週の1位記事" {
		t.Fatalf("weeklyフォールバックが機能していません: %s", pusher.pushed[0].Body)
	}
}

func TestSendDailyDigestNoRanking(t *testing.T) {
	repository := &fakeRepository{tokens: []entity.DeviceToken{newTestToken("ExponentPushToken[aaa]")}}
	service := NewService(repository, &fakeDigestLog{}, &fakePusher{}, &fakeArticleService{articlesByPeriod: map[string][]entity.Article{}})

	if _, err := service.SendDailyDigest(); !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("ランキング全滅時にErrNotFoundが返りません: %v", err)
	}
}

func TestSendDailyDigestSkipsRecentlySent(t *testing.T) {
	first := newTestArticle("1位の記事(送信済み)")
	second := newTestArticle("2位の記事")
	repository := &fakeRepository{tokens: []entity.DeviceToken{newTestToken("ExponentPushToken[aaa]")}}
	pusher := &fakePusher{}
	digestLog := &fakeDigestLog{recentArticleIDs: []string{first.ID.String()}}
	articleService := &fakeArticleService{articlesByPeriod: map[string][]entity.Article{
		"daily": {first, second},
	}}
	service := NewService(repository, digestLog, pusher, articleService)

	if _, err := service.SendDailyDigest(); err != nil {
		t.Fatalf("ダイジェスト送信が失敗しました: %v", err)
	}
	if pusher.pushed[0].Body != "2位の記事" {
		t.Fatalf("送信済み記事が除外されていません: %s", pusher.pushed[0].Body)
	}
	if len(digestLog.savedArticleIDs) != 1 || digestLog.savedArticleIDs[0] != second.ID.String() {
		t.Fatalf("送信履歴が記録されていません: %v", digestLog.savedArticleIDs)
	}
	if digestLog.savedSentCounts[0] != 1 {
		t.Fatalf("送信数の記録が1ではありません: %d", digestLog.savedSentCounts[0])
	}
}

func TestSendDailyDigestAllRecentlySentFallsBackToTop(t *testing.T) {
	first := newTestArticle("1位の記事")
	second := newTestArticle("2位の記事")
	repository := &fakeRepository{tokens: []entity.DeviceToken{newTestToken("ExponentPushToken[aaa]")}}
	pusher := &fakePusher{}
	// 全記事が送信済み → 重複を許容して全体の1位を送る(送らないよりは良い)
	digestLog := &fakeDigestLog{recentArticleIDs: []string{first.ID.String(), second.ID.String()}}
	articleService := &fakeArticleService{articlesByPeriod: map[string][]entity.Article{
		"daily": {first, second},
	}}
	service := NewService(repository, digestLog, pusher, articleService)

	if _, err := service.SendDailyDigest(); err != nil {
		t.Fatalf("ダイジェスト送信が失敗しました: %v", err)
	}
	if pusher.pushed[0].Body != "1位の記事" {
		t.Fatalf("全除外時に1位へフォールバックしていません: %s", pusher.pushed[0].Body)
	}
}

func TestSendDailyDigestNoTokens(t *testing.T) {
	repository := &fakeRepository{}
	pusher := &fakePusher{}
	articleService := &fakeArticleService{articlesByPeriod: map[string][]entity.Article{
		"daily": {newTestArticle("記事")},
	}}
	service := NewService(repository, &fakeDigestLog{}, pusher, articleService)

	sentCount, err := service.SendDailyDigest()
	if err != nil {
		t.Fatalf("トークン0件で失敗しました: %v", err)
	}
	if sentCount != 0 || len(pusher.pushed) != 0 {
		t.Fatal("トークン0件なのに送信されています")
	}
}

func TestSendMatsuri(t *testing.T) {
	repository := &fakeRepository{tokens: []entity.DeviceToken{
		newTestToken("ExponentPushToken[aaa]"),
		newTestToken("ExponentPushToken[bbb]"),
	}}
	pusher := &fakePusher{}
	service := NewService(repository, &fakeDigestLog{}, pusher, &fakeArticleService{})

	matsuriArticle := newTestArticle("祭りの記事")
	sentCount, err := service.SendMatsuri(matsuriArticle, "https://example.com/i.jpg", 4)
	if err != nil {
		t.Fatalf("祭り速報の送信が失敗しました: %v", err)
	}
	if sentCount != 2 {
		t.Fatalf("送信成功数が2ではありません: %d", sentCount)
	}
	if pusher.pushed[0].Title != "🔥 4サイトが一斉にまとめ中" {
		t.Fatalf("通知タイトルが想定と異なります: %s", pusher.pushed[0].Title)
	}
	if pusher.pushed[0].Data["type"] != "matsuri" {
		t.Fatal("通知データにtype=matsuriがありません")
	}
	if pusher.pushed[0].Data["id"] == "" {
		t.Fatal("通知データに記事IDがありません")
	}
}
