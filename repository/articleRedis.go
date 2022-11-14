package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// IncrementViewCount
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

// ListOnlyIDOrderByViewCount
func (r *ArticleRepository) ListOnlyIDOrderByViewCount(period string) ([]entity.ID, error) {
	var idList []entity.ID
	zSetKey, err := getCurrentZSetKey(period)
	if err != nil {
		return nil, err
	}
	serializedMembersWithScores, err := r.kvs.ZRevRangeWithScores(context.Background(), zSetKey, 0, 14).Result()
	for _, serializedMemberWithScore := range serializedMembersWithScores {
		serializedMember := serializedMemberWithScore.Member
		IDString, ok := serializedMember.(string)
		if !ok {
			return nil, err
		}
		ID, err := entity.StringToID(IDString)
		if err != nil {
			return nil, err
		}
		idList = append(idList, ID)
	}
	if len(idList) == 0 {
		return nil, entity.ErrNotFound
	}
	return idList, nil
}

// SetArticleNumber
func (r *ArticleRepository) setArticleNumber(articleNumber int, articleID entity.ID, prefix string) error {
	if err := r.kvs.Set(context.Background(), prefix+":"+strconv.Itoa(articleNumber), articleID.String(), 0).Err(); err != nil {
		return err
	}
	if err := r.kvs.Set(context.Background(), "ml:"+articleID.String(), strconv.Itoa(articleNumber), 0).Err(); err != nil {
		return err
	}
	return nil
}

// GetArticleNumberByArticleID
func (r *ArticleRepository) GetArticleNumberByArticleID(articleID entity.ID, prefix string) (int, error) {
	articleNumberString, err := r.kvs.Get(context.Background(), prefix+":"+articleID.String()).Result()
	if err == redis.Nil {
		return 0, entity.ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(articleNumberString)
}

// GetArticleIDByArticleNumber
func (r *ArticleRepository) getArticleIDByArticleNumber(articleNumber int, prefix string) (entity.ID, error) {
	ArticleIDString, err := r.kvs.Get(context.Background(), prefix+":"+strconv.Itoa(articleNumber)).Result()
	if err != nil {
		return entity.ID{}, err
	}
	return entity.StringToID(ArticleIDString)
}

// getCurrentZSetKey
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
