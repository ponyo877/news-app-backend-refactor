package repository

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic/v7"
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// ArticleElasticSearch elasticsearch repository
type ArticleElasticSearch struct {
	se *elastic.Client
}

type ArticleElasticSearchPresenter struct {
	ID    string
	Title string
}

type ArticleElasticSearchPresenterList []ArticleElasticSearchPresenter

// NewArticleElasticSearch create new repository
func NewArticleElasticSearch(se *elastic.Client) *ArticleElasticSearch {
	return &ArticleElasticSearch{
		se: se,
	}
}

// Search
func (r *ArticleRepository) SearchOnlyID(keyword string) ([]entity.ID, error) {
	var idList []entity.ID
	searchResult, err := r.se.Search().
		Index("article").
		Query(elastic.NewMatchQuery("Title", keyword)).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(searchResult.Hits.Hits) == 0 {
		return []entity.ID{}, nil
	}
	for _, hit := range searchResult.Hits.Hits {
		var articleElasticSearchPresenter ArticleElasticSearchPresenter
		hitSourceByte, err := hit.Source.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(hitSourceByte, &articleElasticSearchPresenter); err != nil {
			return nil, err
		}
		ID, err := entity.StringToID(articleElasticSearchPresenter.ID)
		if err != nil {
			return nil, err
		}
		idList = append(idList, ID)
	}
	return idList, nil
}

// CreateForSearch
func (r *ArticleRepository) CreateForSearch(e entity.Article) error {
	jsonString, err := toJsonString(e)
	if err != nil {
		return err
	}
	if _, err := r.se.Index().
		Index("article").
		Id(e.ID.String()).
		BodyString(jsonString).
		Do(context.Background()); err != nil {
		return err
	}
	return nil
}

// toJsonString
func toJsonString(e entity.Article) (string, error) {
	articleElasticSearchPresenter := &ArticleElasticSearchPresenter{
		ID:    e.ID.String(),
		Title: e.Title.String(),
	}
	jsonData, err := json.Marshal(articleElasticSearchPresenter)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
