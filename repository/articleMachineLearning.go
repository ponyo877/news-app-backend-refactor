package repository

import (
	"context"

	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/pkg/annoyindex"
)

// vectorize
func (r *ArticleRepository) vectorize(title string) ([]float32, error) {
	if r.model == nil {
		return nil, entity.ErrNotFound
	}
	result, err := r.model.Encode(context.Background(), title, int(bert.MeanPooling))
	if err != nil {
		return nil, err
	}
	return result.Vector.Data().F32(), nil
}

// CreateMLIndex
func (r *ArticleRepository) CreateMLIndex(articles []entity.Article) error {
	// ML無効構成(MLM_NAME未設定)では索引を作らない
	if r.model == nil || r.index == nil {
		return entity.ErrNotFound
	}
	newMLIndex := annoyindex.NewAnnoyIndexAngular(256)
	for articleNumber, article := range articles {
		articleTitleVector, err := r.vectorize(article.Title.String())
		if err != nil {
			return err
		}
		if err := r.setArticleNumber(articleNumber, article.ID, "ml"); err != nil {
			return err
		}
		newMLIndex.AddItem(articleNumber, articleTitleVector)
	}
	newMLIndex.Build(10)
	if ok := newMLIndex.Save(r.indexPath); !ok {
		return entity.ErrInternalServerError
	}
	if ok := r.index.Load(r.indexPath); !ok {
		return entity.ErrInternalServerError
	}
	return nil
}

// ListBySimilarity
func (r *ArticleRepository) ListBySimilarity(ID entity.ID) ([]entity.ID, error) {
	// ML無効構成(MLM_NAME未設定)では類似記事なし=空リストとして扱う
	if r.model == nil || r.index == nil {
		return nil, entity.ErrNotFound
	}
	var similarArticleNumbers []int
	var distances []float32
	article, err := r.Get(ID)
	if err != nil {
		return nil, err
	}
	articleTitleVector, err := r.vectorize(article.Title.String())
	if err != nil {
		return nil, err
	}
	r.index.GetNnsByVector(articleTitleVector, 15, -1, &similarArticleNumbers, &distances)
	var idList []entity.ID
	for i, articleNumber := range similarArticleNumbers {
		// 似過ぎている記事は除外
		if distances[i] < 0.1 {
			continue
		}
		articleID, err := r.getArticleIDByArticleNumber(articleNumber, "ml")
		if err != nil {
			return nil, err
		}
		idList = append(idList, articleID)
	}
	return idList, nil
}
