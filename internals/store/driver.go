package store

import (
	"fmt"
	"github.com/alash3al/vecdb/internals/vector"
	"sync"
)

var (
	driversMap      = map[string]Driver{}
	driversMapMutex = &sync.RWMutex{}
)

type Driver interface {
	Open(dsn string) error
	Put(bucket string, key string, vec vector.Vec) error
	Delete(bucket string, key string) error
	Query(q VectorQueryInput) (*VectorQueryResult, error)
	Close() error
}

func Register(name string, driver Driver) {
	driversMapMutex.Lock()
	defer driversMapMutex.Unlock()

	if _, found := driversMap[name]; found {
		panic(fmt.Sprintf("specified driver %s exists", name))
	}

	driversMap[name] = driver
}

func Open(name string, dsn string) (Driver, error) {
	driversMapMutex.RLock()
	defer driversMapMutex.RUnlock()

	driver, found := driversMap[name]
	if !found {
		return nil, fmt.Errorf("driver %s not found", name)
	}

	if err := driver.Open(dsn); err != nil {
		return nil, err
	}

	return driver, nil
}
