//go:build mlindex

package main

import (
	"os"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/pkg/annoyindex"
	"github.com/ponyo877/news-app-backend-refactor/repository"
)

// loadVectorIndex はannoyindex(cgo)を初期化する。mlindexタグ付きビルド専用
func loadVectorIndex(indexPath string) repository.VectorIndex {
	angularIndex := annoyindex.NewAnnoyIndexAngular(256)
	if _, err := os.Stat(indexPath); err == nil {
		if ok := angularIndex.Load(indexPath); ok {
			log.Info("AnnoyIndexの既存モデルの読み込みに成功しました")
		}
	}
	return angularIndex
}
