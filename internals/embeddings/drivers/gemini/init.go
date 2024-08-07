package gemini

import "github.com/alash3al/vecdb/internals/embeddings"

func init() {
	embeddings.Register("gemini", &Driver{})
}
