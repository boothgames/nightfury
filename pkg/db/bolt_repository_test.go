package db_test

import (
	"encoding/json"
	"github.com/boothgames/nightfury/pkg/db"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestModel struct{ Name string }

func (m TestModel) ID() string {
	return m.Name
}

func TestBoltRepositorySave(t *testing.T) {
	dir, _ := ioutil.TempDir("", "nightfury")
	dbPath := path.Join(dir, "db")
	model := TestModel{Name: "should be persisted"}
	repo, _ := db.NewBoltRepository(dbPath)

	defer func() {
		_ = repo.(db.BoltRepository).Close()
		_ = os.RemoveAll(dir)
	}()

	err := repo.Save("test", model)

	assert.NoError(t, err)
}

func TestBoltRepositoryDelete(t *testing.T) {
	dir, _ := ioutil.TempDir("", "nightfury")
	dbPath := path.Join(dir, "db")
	model := TestModel{Name: "should be persisted"}
	repo, _ := db.NewBoltRepository(dbPath)

	defer func() {
		_ = repo.(db.BoltRepository).Close()
		_ = os.RemoveAll(dir)
	}()
	_ = repo.Save("test", model)

	err := repo.Delete("test", model)
	assert.NoError(t, err)

	var actual TestModel
	ok, err := repo.Fetch("test", model.ID(), &actual)

	assert.False(t, ok)
	assert.NoError(t, err)
}

func TestBoltRepositoryFetch(t *testing.T) {
	dir, _ := ioutil.TempDir("", "nightfury")
	dbPath := path.Join(dir, "db")
	repo, _ := db.NewBoltRepository(dbPath)
	model := TestModel{Name: "should be persisted"}

	defer func() {
		_ = repo.(db.BoltRepository).Close()
		_ = os.RemoveAll(dir)
	}()

	err := repo.Save("test", model)
	assert.NoError(t, err)

	var actual TestModel
	_, err = repo.Fetch("test", "should be persisted", &actual)
	assert.NoError(t, err)

	if !cmp.Equal(model, actual) {
		assert.Fail(t, cmp.Diff(model, actual))
	}
}

func TestBoltRepositoryFetchAll(t *testing.T) {
	dir, _ := ioutil.TempDir("", "nightfury")
	dbPath := path.Join(dir, "db")
	repo, _ := db.NewBoltRepository(dbPath)
	modelOne := TestModel{Name: "one"}
	modelTwo := TestModel{Name: "two"}
	modelThree := TestModel{Name: "three"}

	defer func() {
		_ = repo.(db.BoltRepository).Close()
		_ = os.RemoveAll(dir)
	}()

	_ = repo.Save("test", modelOne)
	_ = repo.Save("test", modelTwo)
	_ = repo.Save("test", modelThree)

	expected := map[string]interface{}{"one": modelOne, "two": modelTwo, "three": modelThree}

	actual, err := repo.FetchAll("test", func(bytes []byte) (db.Model, error) {
		model := TestModel{}
		err := json.Unmarshal(bytes, &model)
		return model, err
	})

	assert.NoError(t, err)

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}
