package db_test

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestModel string

func (m TestModel) ID() string {
	return fmt.Sprintf("%v", m)
}

func TestBoltRepositorySave(t *testing.T) {
	dir, _ := ioutil.TempDir("", "nightfury")
	dbPath := path.Join(dir, "db")
	model := TestModel("should be persisted")
	repo, _ := db.NewBoltRepository(dbPath)

	defer func() {
		_ = repo.(db.BoltRepository).Close()
		_ = os.RemoveAll(dir)
	}()

	err := repo.Save("test", model)

	assert.NoError(t, err)
}

func TestBoltRepositoryFetch(t *testing.T) {
	dir, _ := ioutil.TempDir("", "nightfury")
	dbPath := path.Join(dir, "db")
	repo, _ := db.NewBoltRepository(dbPath)
	model := TestModel("should be persisted")

	defer func() {
		_ = repo.(db.BoltRepository).Close()
		_ = os.RemoveAll(dir)
	}()

	err := repo.Save("test", model)
	assert.NoError(t, err)

	var actual TestModel
	err = repo.Fetch("test", "should be persisted", &actual)
	assert.NoError(t, err)

	if !cmp.Equal(model, actual) {
		assert.Fail(t, cmp.Diff(model, actual))
	}
}
