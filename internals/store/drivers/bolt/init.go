package bolt

import "github.com/alash3al/vecdb/internals/store"

func init() {
	store.Register("bolt", &Driver{})
}
