package notification

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
)

// Service Service struct
type Service struct {
	repository     Repository
	pusher         Pusher
	articleService article.UseCase
}

// NewService create new service
func NewService(r Repository, p Pusher, a article.UseCase) *Service {
	return &Service{
		repository:     r,
		pusher:         p,
		articleService: a,
	}
}

// RegisterToken トークンの登録・設定更新(同一トークンはupsert)
func (s *Service) RegisterToken(expoToken, deviceHash, platform string, digestEnabled bool) error {
	deviceToken, err := entity.NewDeviceToken(expoToken, deviceHash, platform, digestEnabled)
	if err != nil {
		return err
	}
	return s.repository.Save(deviceToken)
}

// SendDailyDigest 人気1位の記事をダイジェスト通知として全許諾端末へ送る。
// 朝の実行時はdailyランキングがまだ薄いため weekly → monthly へフォールバックする
func (s *Service) SendDailyDigest() (int, error) {
	topArticle, err := s.pickTopArticle()
	if err != nil {
		return 0, err
	}
	deviceTokenList, err := s.repository.ListDigestEnabled()
	if err != nil {
		return 0, err
	}
	if len(deviceTokenList) == 0 {
		return 0, nil
	}
	// タップ時にアプリが記事画面を直接開けるよう、一覧APIと同じメタ情報一式を積む
	imageURL, err := topArticle.ImageURL.URL()
	if err != nil {
		imageURL = ""
	}
	data := map[string]string{
		"type":        "digest",
		"id":          topArticle.ID.String(),
		"titles":      topArticle.Title.String(),
		"url":         topArticle.URL,
		"image":       imageURL,
		"siteID":      topArticle.Site.ID.String(),
		"sitetitle":   topArticle.Site.Title,
		"publishedAt": topArticle.PublishedAt.Format(time.RFC3339),
	}
	messages := make([]entity.PushMessage, 0, len(deviceTokenList))
	for _, deviceToken := range deviceTokenList {
		messages = append(messages, entity.PushMessage{
			To:    deviceToken.ExpoToken,
			Title: "今日の人気No.1",
			Body:  topArticle.Title.String(),
			Data:  data,
		})
	}
	invalidTokens, err := s.pusher.Push(messages)
	if err != nil {
		return 0, err
	}
	for _, invalidToken := range invalidTokens {
		if err := s.repository.Delete(invalidToken); err != nil {
			log.Warnf("失効トークンの削除に失敗しました(%s): %v", invalidToken, err)
		}
	}
	return len(messages) - len(invalidTokens), nil
}

func (s *Service) pickTopArticle() (entity.Article, error) {
	var lastErr error
	for _, period := range []string{"daily", "weekly", "monthly"} {
		articles, err := s.articleService.ListPopularArticles(period)
		if err != nil {
			lastErr = err
			continue
		}
		if len(articles) > 0 {
			return articles[0], nil
		}
	}
	return entity.Article{}, lastErr
}
