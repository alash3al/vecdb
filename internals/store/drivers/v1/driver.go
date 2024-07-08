package v1

import (
	"encoding/json"
	"errors"
	"github.com/alash3al/vecdb/internals/store"
	"github.com/alash3al/vecdb/internals/vector"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
)

var _ store.Driver = (*Driver)(nil)
var errStop = errors.New("STOP")

type Driver struct {
	db     *bolt.DB
	bucket *bolt.Bucket
}

func (d *Driver) Open(dsn string) error {
	db, err := bolt.Open(dsn, 0600, bolt.DefaultOptions)
	if err != nil {
		return err
	}

	d.db = db

	return nil
}

func (d *Driver) Put(bucket string, key string, vec vector.Vec) error {
	keyBytes := []byte(key)
	valueBytes, err := json.Marshal(vec)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return bucket.Put(keyBytes, valueBytes)
	})
}

func (d *Driver) Delete(bucket string, key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return bucket.Delete([]byte(key))
	})
}

func (d *Driver) Query(q store.VectorQueryInput) (*store.VectorQueryResult, error) {
	var err error

	result := new(store.VectorQueryResult)

	err = d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(q.Bucket))

		return err
	})

	if err != nil {
		return nil, err
	}

	err = d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(q.Bucket))

		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		return bucket.ForEach(func(k, v []byte) error {
			var currentIterationVec vector.Vec

			if err := json.Unmarshal(v, &currentIterationVec); err != nil {
				return err
			}

			if distance := currentIterationVec.CosineSimilarity(q.Vector); distance >= q.MinCosineSimilarityDistance {
				result.Items = append(result.Items, store.VectorQueryResultItem{
					Key:      string(k),
					Distance: distance,
				})
			}

			result.Count++

			if result.Count >= q.MaxResultCount {
				return errStop
			}

			return nil
		})
	})

	if err != nil && !errors.Is(err, errStop) {
		return nil, err
	}

	slices.SortFunc[store.VectorQueryResultItem](result.Items, func(a, b store.VectorQueryResultItem) bool {
		return a.Distance > b.Distance
	})

	return result, nil
}

func (d *Driver) Close() error {
	return d.db.Close()
}
