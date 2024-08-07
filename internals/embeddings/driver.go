package embeddings

import (
	"context"
	"fmt"
	"github.com/alash3al/vecdb/internals/vector"
	"sync"
)

var (
	driversMap      = map[string]Driver{}
	driversMapMutex = &sync.RWMutex{}
)

type Driver interface {
	Open(args map[string]any) error
	TextEmbedding(ctx context.Context, content string) (vector.Vec, error)
	Close() error
}

func Register(name string, driver Driver) {
	driversMapMutex.Lock()
	defer driversMapMutex.Unlock()

	if _, found := driversMap[name]; found {
		panic(fmt.Sprintf("specified embeddings driver %s exists", name))
	}

	driversMap[name] = driver
}

func Open(name string, args map[string]any) (Driver, error) {
	driversMapMutex.RLock()
	defer driversMapMutex.RUnlock()

	driver, found := driversMap[name]
	if !found {
		return nil, fmt.Errorf("embeddings driver %s not found", name)
	}

	if err := driver.Open(args); err != nil {
		return nil, err
	}

	return driver, nil
}
