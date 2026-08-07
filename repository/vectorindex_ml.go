//go:build mlindex

package repository

import "github.com/ponyo877/news-app-backend-refactor/pkg/annoyindex"

// mlindexタグ付きビルドでのみannoyindex(cgo)を配線する
func init() {
	newVectorIndex = func() VectorIndex {
		return annoyindex.NewAnnoyIndexAngular(256)
	}
}
