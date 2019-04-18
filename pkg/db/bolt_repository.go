package db

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

// BoltRepository represents a bbolt database
type BoltRepository struct {
	db *bbolt.DB
}

// NewBoltRepository returns the repository and error if any
func NewBoltRepository(path string) (Repository, error) {
	instance, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open db, reason %v", err)
	}
	return BoltRepository{db: instance}, nil
}

// Close closes bbolt db connection
func (b BoltRepository) Close() error {
	return b.db.Close()
}

// Save persists the model in the bucketName
func (b BoltRepository) Save(bucketName string, model Model) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bName := []byte(bucketName)
		bucket := tx.Bucket(bName)
		if bucket == nil {
			clientsBucket, err := tx.CreateBucket(bName)
			if err != nil {
				return err
			}
			bucket = clientsBucket
		}

		bytes, err := json.Marshal(model)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(model.ID()), bytes)
	})
}

// Fetch retrieves the model identified by name from the bucket named bucketName
func (b BoltRepository) Fetch(bucketName string, name string, model Model) error {
	return b.db.View(func(tx *bbolt.Tx) error {
		bName := []byte(bucketName)
		bucket := tx.Bucket(bName)
		if bucket == nil {
			return nil
		}
		if data := bucket.Get([]byte(name)); data != nil {
			return json.Unmarshal(data, model)
		}
		return nil

	})
}

// FetchAll returns all the models available in the bucketName and error if any
func (b BoltRepository) FetchAll(bucketName string, modelFn func([]byte) (Model, error)) (interface{}, error) {
	result := map[string]interface{}{}
	err := b.db.View(func(tx *bbolt.Tx) error {
		bName := []byte(bucketName)
		bucket := tx.Bucket(bName)
		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(key, value []byte) error {
			model, err := modelFn(value)
			if err == nil {
				result[string(key)] = model
			}
			return err
		})
		return err
	})
	return result, err
}
