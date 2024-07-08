package store

import (
	"github.com/alash3al/vecdb/internals/vector"
)

type VectorQueryInput struct {
	Bucket                      string
	Vector                      vector.Vec
	MinCosineSimilarityDistance float64
	MaxResultCount              int64
}

type VectorQueryResult struct {
	Count int64
	Items []VectorQueryResultItem
}

type VectorQueryResultItem struct {
	Key      string
	Distance float64
}
