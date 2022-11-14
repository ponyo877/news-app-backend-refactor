package repository

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"github.com/ponyo877/news-app-backend-refactor/pkg/annoyindex"
)

type Vector = []float32

// Vectorize
func (r *ArticleRepository) vectorize(title string) (Vector, error) {
	result, err := r.model.Encode(context.Background(), title, int(bert.MeanPooling))
	if err != nil {
		return nil, err
	}
	return result.Vector.Data().F32(), nil
}

// AddArticles
func (r *ArticleRepository) CreateMLIndex(articles []entity.Article) error {
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
	var similarArticleNumbers []int
	targetArticleNumber, err := r.GetArticleNumberByArticleID(ID, "ml")
	if err != nil && err != entity.ErrNotFound {
		return nil, err
	}
	// まだMLIndexに登録されてない場合はタイトルベクトルから計算する
	if err == entity.ErrNotFound {
		article, err := r.Get(ID)
		if err != nil {
			return nil, err
		}
		articleTitleVector, err := r.vectorize(article.Title.String())
		if err != nil {
			return nil, err
		}
		r.index.GetNnsByVector(articleTitleVector, 15, -1, &similarArticleNumbers)
	} else {
		log.Infof("articleID: %v, targetArticleNumber: %v", ID.String(), targetArticleNumber)
		r.index.GetNnsByItem(targetArticleNumber, 15, -1, &similarArticleNumbers)
	}

	var idList []entity.ID
	for _, articleNumber := range similarArticleNumbers {
		articleID, err := r.getArticleIDByArticleNumber(articleNumber, "ml")
		if err != nil {
			return nil, err
		}
		idList = append(idList, articleID)
	}
	return idList, nil
}
