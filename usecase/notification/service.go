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
	digestLog      DigestLogRepository
	pusher         Pusher
	articleService article.UseCase
}

// NewService create new service
func NewService(r Repository, d DigestLogRepository, p Pusher, a article.UseCase) *Service {
	return &Service{
		repository:     r,
		digestLog:      d,
		pusher:         p,
		articleService: a,
	}
}

// 直近この時間内に送った記事は次回の選定から除外する。
// 24時間だと朝夜の2枠しか外れないが、48時間なら「1位に数日居座る記事」の再送も防げる
const digestDedupWindow = 48 * time.Hour

// RegisterToken トークンの登録・設定更新(同一トークンはupsert)
func (s *Service) RegisterToken(expoToken, deviceHash, platform string, digestEnabled bool) error {
	deviceToken, err := entity.NewDeviceToken(expoToken, deviceHash, platform, digestEnabled)
	if err != nil {
		return err
	}
	return s.repository.Save(deviceToken)
}

// SendDailyDigest 人気1位の記事をダイジェスト通知として全許諾端末へ送る。
// 朝の実行時はdailyランキングがまだ薄いため weekly → monthly へフォールバックする。
// 直近に送った記事は除外する(朝と夜で同じ記事が2回届くのを防ぐ)
func (s *Service) SendDailyDigest() (int, error) {
	recentArticleIDs, err := s.digestLog.ListRecentArticleIDs(time.Now().Add(-digestDedupWindow))
	if err != nil {
		// 履歴が読めなくても送信自体は止めない(最悪でも重複するだけ)
		log.Warnf("ダイジェスト送信履歴の取得に失敗しました: %v", err)
		recentArticleIDs = nil
	}
	exclude := make(map[string]bool, len(recentArticleIDs))
	for _, articleID := range recentArticleIDs {
		exclude[articleID] = true
	}
	topArticle, err := s.pickTopArticle(exclude)
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
	sentCount := len(messages) - len(invalidTokens)
	if err := s.digestLog.Save(topArticle.ID, sentCount); err != nil {
		log.Warnf("ダイジェスト送信履歴の記録に失敗しました: %v", err)
	}
	return sentCount, nil
}

// pickTopArticle 除外リストにない最上位の記事を選ぶ。
// 全記事が除外済みの場合のみ重複を許容して全体の1位を返す(送らないよりは良い)
func (s *Service) pickTopArticle(exclude map[string]bool) (entity.Article, error) {
	var fallback *entity.Article
	var lastErr error
	for _, period := range []string{"daily", "weekly", "monthly"} {
		articles, err := s.articleService.ListPopularArticles(period)
		if err != nil {
			lastErr = err
			continue
		}
		for i := range articles {
			if fallback == nil {
				fallback = &articles[i]
			}
			if !exclude[articles[i].ID.String()] {
				return articles[i], nil
			}
		}
	}
	if fallback != nil {
		return *fallback, nil
	}
	return entity.Article{}, lastErr
}
