package entity

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
	return false
}

// Validate
func (a ArticleTitle) Validate() error {
	return nil
}
