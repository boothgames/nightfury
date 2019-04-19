package db

import "go.etcd.io/bbolt"

var boltRepository BoltRepository

// Model which can persisted in the repository
type Model interface {
	ID() string
}

// Repository holds the necessary method to persist and retrieve data
// from database
type Repository interface {
	Save(bucketName string, model Model) error
	Delete(bucketName string, model Model) error
	Fetch(bucketName string, name string, model Model) (bool, error)
	FetchAll(bucketName string, modelFn func(data []byte) (Model, error)) (interface{}, error)
}

// Initialize initializes the global repository
func Initialize(path string) error {
	repository, err := NewBoltRepository(path)
	if err != nil {
		return err
	}
	boltRepository = repository.(BoltRepository)
	return nil
}

// DeleteBucket will delete the bucket from the db
func DeleteBucket(bucketName string) error {
	err := boltRepository.db.Update(func(tx *bbolt.Tx) error {
		bName := []byte(bucketName)
		bucket := tx.Bucket(bName)
		if bucket != nil {
			err := tx.DeleteBucket(bName)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Close the default repository
func Close() error {
	return boltRepository.Close()
}

// DefaultRepository returns the global repository
func DefaultRepository() Repository {
	return boltRepository
}
