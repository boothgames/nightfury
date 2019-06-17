package nightfury_test

import (
	"encoding/json"
	"fmt"
	"github.com/boothgames/nightfury/pkg/db"
	mocks "github.com/boothgames/nightfury/pkg/internal/mocks/db"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHintID(t *testing.T) {
	t.Run("should return the id", func(t *testing.T) {
		hint := nightfury.Hint{
			Title:    "sample hInT",
			Content:  "content",
			Tag:      []string{"web"},
			Takeaway: "dont do this",
		}

		actual := hint.ID()

		assert.Equal(t, "sample-hint", actual)
	})
}

func TestHintSave(t *testing.T) {
	t.Run("should be able to save hint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		hint := nightfury.Hint{
			Title:    "title",
			Content:  "content",
			Tag:      []string{"web"},
			Takeaway: "dont do this",
		}
		repository.EXPECT().Save("hints", hint)

		err := hint.Save(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		hint := nightfury.Hint{Title: "title"}
		repository.EXPECT().Save("hints", hint).Return(fmt.Errorf("unable to save"))

		err := hint.Save(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestNewHintFromRepoWithName(t *testing.T) {
	t.Run("should fetch the hint from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("hints", "one", gomock.Any()).Return(true, nil)

		actual, err := nightfury.NewHintFromRepoWithName(repository, "one")

		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("should fail to fetch the client from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("hints", "one", gomock.Any()).Return(false, nil)

		actual, err := nightfury.NewHintFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "hint with name one doesn't exists", err.Error())
		}
		assert.Equal(t, nightfury.Hint{}, actual)
	})

	t.Run("should return error returned while fetching data from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("hints", "one", gomock.Any()).Return(false, fmt.Errorf("unable to fetch"))

		actual, err := nightfury.NewHintFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch", err.Error())
		}
		assert.Equal(t, nightfury.Hint{}, actual)
	})
}

func TestNewHintsFromRepo(t *testing.T) {
	t.Run("should be able to get all the hints", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := nightfury.Hints{
			"title": nightfury.Hint{
				Title:    "title",
				Content:  "content",
				Tag:      []string{"web"},
				Takeaway: "dont do this",
			},
		}
		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().FetchAll("hints", gomock.Any()).DoAndReturn(
			func(bucketName string, modelFn func(data []byte) (db.Model, error)) (interface{}, error) {
				data, _ := json.Marshal(nightfury.Hint{
					Title:    "title",
					Content:  "content",
					Tag:      []string{"web"},
					Takeaway: "dont do this",
				})
				model, err := modelFn(data)
				if err != nil {
					return nil, err
				}
				return nightfury.Hints{model.ID(): model.(nightfury.Hint)}, nil
			})

		hints, err := nightfury.NewHintsFromRepo(repository)

		assert.NoError(t, err)
		if !cmp.Equal(expected, hints) {
			assert.Fail(t, cmp.Diff(nightfury.Hint{}, hints))
		}
	})
}

func TestHintDelete(t *testing.T) {
	t.Run("should be able to delete hint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		hint := nightfury.Hint{Title: "hint"}
		repository.EXPECT().Delete("hints", hint)

		err := hint.Delete(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		hint := nightfury.Hint{Title: "hint"}
		repository.EXPECT().Delete("hints", hint).Return(fmt.Errorf("unable to delete"))

		err := hint.Delete(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to delete", err.Error())
		}
	})
}
