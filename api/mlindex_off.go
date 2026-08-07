//go:build !mlindex

package main

import (
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/repository"
)

// loadVectorIndex は通常ビルド(cgoなし)では類似記事インデックスを提供しない。
// クロスコンパイル(CGO_ENABLED=0)を可能にするための空実装
func loadVectorIndex(_ string) repository.VectorIndex {
	log.Info("mlindexタグなしビルドのため類似記事インデックスは無効です")
	return nil
}
