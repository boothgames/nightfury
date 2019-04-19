package db

var repository Repository

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
	Close() error
}

// Initialize initializes the global repository
func Initialize(path string) error {
	repo, err := NewBoltRepository(path)
	if err != nil {
		return err
	}
	repository = repo.(BoltRepository)
	return nil
}

// Close the default repository
func Close() error {
	return repository.Close()
}

// DefaultRepository returns the global repository
func DefaultRepository() Repository {
	return repository
}

// ReplaceDefaultRepositoryWith replace the default repository
func ReplaceDefaultRepositoryWith(repo Repository) func() {
	originalRepo := repository
	repository = repo
	return func() {
		repository = originalRepo
	}
}
