package entity

import (
	"regexp"

	"github.com/tomoemon/text_normalizer"
)

type ArticleTitle struct {
	Title string
}

// NewArticleTitle create a new entity ArticleTitle
func NewArticleTitle(title string) ArticleTitle {
	return ArticleTitle{
		Title: title,
	}
}

// String
func (a ArticleTitle) String() string {
	return a.Title
}

// ContainsBlacklist
func (a ArticleTitle) ContainsBlacklist() bool {
	return NewBlackList(a.normalizedTitle()).Contain()
}

// Validate
func (a ArticleTitle) Validate() error {
	return nil
}

// normalize
func (a ArticleTitle) normalizedTitle() string {
	normalizer := text_normalizer.NewTextNormalizer(
		text_normalizer.AlphabetToHankaku,
		text_normalizer.KanaToHiragana,
	)
	rep := regexp.MustCompile("[^0-9a-zA-Zぁ-んァ-ヶ一-龠ー]+")
	return rep.ReplaceAllString(normalizer.Replace(a.Title), "@")
}
