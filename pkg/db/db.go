package db

var boltRepository BoltRepository

// Model which can persisted in the repository
type Model interface {
	ID() string
}

// Repository holds the necessary method to persist and retrieve data
// from database
type Repository interface {
	Save(bucketName string, model Model) error
	Fetch(bucketName string, name string, model Model) error
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

// Close the default repository
func Close() error {
	return boltRepository.Close()
}

// DefaultRepository returns the global repository
func DefaultRepository() Repository {
	return boltRepository
}
