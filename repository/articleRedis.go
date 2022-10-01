package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// ArticleRedis redis repository
type ArticleRedis struct {
	kvs *redis.Client
}

type ArticleRedisPresenter struct {
	ID        string
	ViewCount int
	Rank      int
}

type ArticleRedisPresenterList []ArticleRedisPresenter

// NewArticleRedis create new repository
func NewArticleRedis(kvs *redis.Client) *ArticleRedis {
	return &ArticleRedis{
		kvs: kvs,
	}
}

func (r *ArticleRepository) IncrementViewCount(ID entity.ID) error {
	periodList := []string{"monthly", "weekly", "daily"}
	for _, period := range periodList {
		zSetKey, err := getCurrentZSetKey(period)
		if err != nil {
			return err
		}
		if _, err := r.kvs.ZIncrBy(context.Background(), zSetKey, 1, ID.String()).Result(); err != nil {
			return err
		}
	}
	return nil
}

func (r *ArticleRepository) ListOnlyIDOrderByViewCount(period string) ([]entity.ID, error) {
	// var articleRedisPresenterList ArticleRedisPresenterList
	var idList []entity.ID
	zSetKey, err := getCurrentZSetKey(period)
	if err != nil {
		return nil, err
	}
	serializedMembersWithScores, err := r.kvs.ZRevRangeWithScores(context.Background(), zSetKey, 0, 14).Result()
	// member := &Member{}
	for _, serializedMemberWithScore := range serializedMembersWithScores {
		serializedMember := serializedMemberWithScore.Member
		// score := serializedMemberWithScore.Score
		// if err := json.Unmarshal([]byte(serializedMember.(string)), member); err != nil {
		// 	return nil, err
		// }
		ID, ok := serializedMember.(string)
		if !ok {
			return nil, err
		}
		id, err := entity.StringToID(ID)
		if err != nil {
			return nil, err
		}
		idList = append(idList, id)
		// articleRedisPresenter := ArticleRedisPresenter{
		// 	ID:        member.ID,
		// 	ViewCount: int(score),
		// 	Rank:      idx + 1,
		// }
		// articleRedisPresenterList = append(articleRedisPresenterList, articleRedisPresenter)
	}
	return idList, nil
}

func getCurrentZSetKey(period string) (string, error) {
	periodToZSetKey := map[string]string{}
	now := time.Now()
	lastSunday := now.AddDate(0, 0, int(time.Sunday-now.Weekday()))

	periodToZSetKey["monthly"] = fmt.Sprintf("m_%s-01", now.Format("2006-01"))
	periodToZSetKey["weekly"] = fmt.Sprintf("w_%s", lastSunday.Format("2006-01-02"))
	periodToZSetKey["daily"] = fmt.Sprintf("d_%s", now.Format("2006-01-02"))
	zSetKey, ok := periodToZSetKey[period]
	if !ok {
		return "", entity.ErrNotFound
	}
	return zSetKey, nil
}
