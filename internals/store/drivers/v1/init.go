package v1

import "github.com/alash3al/vecdb/internals/store"

func init() {
	store.Register("v1", &Driver{})
}
