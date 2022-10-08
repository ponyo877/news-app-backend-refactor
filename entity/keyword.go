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

// QueryArg
func (k *Keyword) QueryArg() string {
	keywordList := strings.FieldsFunc(k.Value, isSpace)
	queryArg := strings.Join(keywordList, " ")
	return queryArg
}

// isSpace
func isSpace(r rune) bool {
	return r == ' ' || r == '　'
}
