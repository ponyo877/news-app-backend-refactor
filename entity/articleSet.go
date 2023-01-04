package entity

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type ArticleSiteAndTitle struct {
	ID    ID
	Title ArticleTitle
}

type ArticleSet struct {
	Set      []Article
	CheckSet mapset.Set[ArticleSiteAndTitle]
}

// NewArticleSet create a new entity ArticleSet
func NewArticleSet() ArticleSet {
	return ArticleSet{
		Set:      []Article{},
		CheckSet: mapset.NewSet[ArticleSiteAndTitle](),
	}
}

// Add
func (a ArticleSet) Add(article Article) ArticleSet {
	articleSiteAndTitle := ArticleSiteAndTitle{
		ID:    article.ID,
		Title: article.Title,
	}
	if a.CheckSet.Contains(articleSiteAndTitle) {
		return a
	}
	a.CheckSet.Add(articleSiteAndTitle)
	a.Set = append(a.Set, article)
	return a
}
