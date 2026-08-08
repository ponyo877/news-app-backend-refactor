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

func (f *fakeRepository) Delete(expoToken string) error {
	f.deleted = append(f.deleted, expoToken)
	return nil
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
	service := NewService(repository, &fakePusher{}, &fakeArticleService{})

	if err := service.RegisterToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]", "devicehash", "ios", true); err != nil {
		t.Fatalf("正常なトークン登録が失敗しました: %v", err)
	}
	if len(repository.saved) != 1 {
		t.Fatalf("保存件数が1ではありません: %d", len(repository.saved))
	}

	if err := service.RegisterToken("bogus-token", "devicehash", "ios", true); err == nil {
		t.Fatal("不正な形式のトークンが登録できてしまいました")
	}
	if err := service.RegisterToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]", "devicehash", "web", true); err == nil {
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
	service := NewService(repository, pusher, articleService)

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
	service := NewService(repository, pusher, articleService)

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
	service := NewService(repository, &fakePusher{}, &fakeArticleService{articlesByPeriod: map[string][]entity.Article{}})

	if _, err := service.SendDailyDigest(); !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("ランキング全滅時にErrNotFoundが返りません: %v", err)
	}
}

func TestSendDailyDigestNoTokens(t *testing.T) {
	repository := &fakeRepository{}
	pusher := &fakePusher{}
	articleService := &fakeArticleService{articlesByPeriod: map[string][]entity.Article{
		"daily": {newTestArticle("記事")},
	}}
	service := NewService(repository, pusher, articleService)

	sentCount, err := service.SendDailyDigest()
	if err != nil {
		t.Fatalf("トークン0件で失敗しました: %v", err)
	}
	if sentCount != 0 || len(pusher.pushed) != 0 {
		t.Fatal("トークン0件なのに送信されています")
	}
}
