package entity

import (
	"strings"
)

type Keyword struct {
	Value string
}

// NewComment create a new article
func NewKeyword(keyword string) (Keyword, error) {
	newKeyword := Keyword{
		Value: keyword,
	}
	if err := newKeyword.Validate(); err != nil {
		return Keyword{}, ErrInvalidEntity
	}
	return newKeyword, nil
}

// Validate validate data
func (k *Keyword) Validate() error {
	if k.Value == "" {
		return ErrInvalidEntity
	}
	return nil
}

// QueryArg はBOOLEAN MODE検索の引数を組み立てる。
// 各語に+を付けたAND検索(スペース区切り=絞り込みという一般的な検索UXに合わせる)。
// ユーザー入力に含まれる演算子記号はクエリ構文エラーの原因になるため除去する。
func (k *Keyword) QueryArg() string {
	keywordList := strings.FieldsFunc(k.Value, isSpace)
	queryArgList := make([]string, 0, len(keywordList))
	for _, word := range keywordList {
		word = strings.Map(dropBooleanOperator, word)
		if word == "" {
			continue
		}
		queryArgList = append(queryArgList, "+"+word)
	}
	return strings.Join(queryArgList, " ")
}

// dropBooleanOperator はMySQL BOOLEAN MODEで演算子として解釈される記号を落とす
func dropBooleanOperator(r rune) rune {
	switch r {
	case '+', '-', '<', '>', '(', ')', '~', '*', '"', '@':
		return -1
	}
	return r
}

// isSpace
func isSpace(r rune) bool {
	return r == ' ' || r == '　'
}
