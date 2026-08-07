package repository

// VectorIndex は近似最近傍インデックス(annoyindex)の抽象。
// annoyindexはcgo(C++)依存でクロスコンパイルを阻害するため、
// ビルドタグ mlindex の背後に隔離する境界としてこのインターフェースを置く。
// シグネチャはSWIG生成コード(pkg/annoyindex)のAnnoyIndexに合わせている。
type VectorIndex interface {
	AddItem(item int, vector []float32)
	Build(nTrees int)
	Save(a ...interface{}) bool
	Load(a ...interface{}) bool
	GetNnsByVector(a ...interface{})
}

// newVectorIndex は mlindex ビルド時のみ vectorindex_ml.go のinitで設定される。
// nil のままなら類似記事インデックスの新規構築は無効(通常ビルド)
var newVectorIndex func() VectorIndex
