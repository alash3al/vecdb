package store

import (
	"github.com/alash3al/vecdb/internals/vector"
)

type VectorQueryInput struct {
	Bucket              string
	Vector              vector.Vec
	MinCosineSimilarity float64
	MaxResultCount      int64
}

type VectorQueryResult struct {
	Items []VectorQueryResultItem `json:"items" json:"items,omitempty"`
}

type VectorQueryResultItem struct {
	Key              string  `json:"key"`
	CosineSimilarity float64 `json:"cosine_similarity"`
}
