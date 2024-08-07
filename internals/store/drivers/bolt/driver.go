package bolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alash3al/vecdb/internals/store"
	"github.com/alash3al/vecdb/internals/vector"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
	"math"
)

var _ store.Driver = (*Driver)(nil)
var errStop = errors.New("STOP")

type Driver struct {
	db     *bolt.DB
	bucket *bolt.Bucket
}

type Value struct {
	Vector       vector.Vec
	Magnitude    float64
	SigmoidValue float64
}

func (d *Driver) Open(args map[string]any) error {
	filename, ok := args["database"].(string)
	if !ok {
		return fmt.Errorf("unable to find the store `args.database`")
	}

	db, err := bolt.Open(filename, 0600, bolt.DefaultOptions)
	if err != nil {
		return err
	}

	d.db = db

	return nil
}

func (d *Driver) Put(bucket string, key string, vec vector.Vec) error {
	keyBytes := []byte(key)
	valueBytes, err := json.Marshal(Value{
		Vector:       vec,
		Magnitude:    vec.Magnitude(),
		SigmoidValue: 1 / (1 + math.Exp(-1*vec.Magnitude())),
	})
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
			var currentIterationValue Value

			if err := json.Unmarshal(v, &currentIterationValue); err != nil {
				return err
			}

			if cosineSimilarity := currentIterationValue.Vector.CosineSimilarity(q.Vector); cosineSimilarity >= q.MinCosineSimilarity {
				result.Items = append(result.Items, store.VectorQueryResultItem{
					Key:              string(k),
					CosineSimilarity: cosineSimilarity,
				})
			}

			return nil
		})
	})

	if err != nil && !errors.Is(err, errStop) {
		return nil, err
	}

	slices.SortFunc[store.VectorQueryResultItem](result.Items, func(a, b store.VectorQueryResultItem) bool {
		return a.CosineSimilarity > b.CosineSimilarity
	})

	if q.MaxResultCount > 0 {
		result.Items = result.Items[0:q.MaxResultCount]
	}

	return result, nil
}

func (d *Driver) Close() error {
	return d.db.Close()
}
