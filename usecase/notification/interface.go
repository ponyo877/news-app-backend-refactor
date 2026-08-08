package notification

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// Repository interface
type Repository interface {
	Save(e entity.DeviceToken) error
	ListDigestEnabled() ([]entity.DeviceToken, error)
	Delete(expoToken string) error
}

// DigestLogRepository interface(送信履歴: 重複送信防止とCTR計測の分母)
type DigestLogRepository interface {
	Save(articleID entity.ID, sentCount int) error
	ListRecentArticleIDs(since time.Time) ([]string, error)
}

// Pusher interface(Expo Push Service)
type Pusher interface {
	Push(messages []entity.PushMessage) (invalidTokens []string, err error)
}

// UseCase interface
type UseCase interface {
	RegisterToken(expoToken, deviceHash, platform string, digestEnabled bool) error
	SendDailyDigest() (int, error)
}
