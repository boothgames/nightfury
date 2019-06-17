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

func TestGameID(t *testing.T) {
	t.Run("should return the id", func(t *testing.T) {
		game := nightfury.Game{
			Name: "Sample_gAme",
		}

		actual := game.ID()

		assert.Equal(t, "sample-game", actual)
	})
}

func TestGameSave(t *testing.T) {
	t.Run("should be able to save game", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		game := nightfury.Game{Name: "game"}
		repository.EXPECT().Save("games", game)

		err := game.Save(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		game := nightfury.Game{Name: "game"}
		repository.EXPECT().Save("games", game).Return(fmt.Errorf("unable to save"))

		err := game.Save(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestGameDelete(t *testing.T) {
	t.Run("should be able to save game", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		game := nightfury.Game{Name: "game"}
		repository.EXPECT().Delete("games", game)

		err := game.Delete(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		game := nightfury.Game{Name: "game"}
		repository.EXPECT().Delete("games", game).Return(fmt.Errorf("unable to save"))

		err := game.Delete(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestNewGameFromRepoWithName(t *testing.T) {
	t.Run("should fetch the game from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("games", "one", gomock.Any()).Return(true, nil)

		actual, err := nightfury.NewGameFromRepoWithName(repository, "one")

		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("should fail to fetch the client from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("games", "one", gomock.Any()).Return(false, nil)

		actual, err := nightfury.NewGameFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "game with name one doesn't exists", err.Error())
		}
		assert.Equal(t, nightfury.Game{}, actual)
	})

	t.Run("should return error returned while fetching data from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("games", "one", gomock.Any()).Return(false, fmt.Errorf("unable to fetch"))

		actual, err := nightfury.NewGameFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch", err.Error())
		}
		assert.Equal(t, nightfury.Game{}, actual)
	})
}

func TestNewGamesFromRepo(t *testing.T) {
	t.Run("should be able to get all the games", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := nightfury.Games{"example": {Name: "example", Instruction: "instruction"}}
		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().FetchAll("games", gomock.Any()).DoAndReturn(
			func(bucketName string, modelFn func(data []byte) (db.Model, error)) (interface{}, error) {
				data, _ := json.Marshal(nightfury.Game{Name: "example", Instruction: "instruction"})
				model, err := modelFn(data)
				if err != nil {
					return nil, err
				}
				return nightfury.Games{model.ID(): model.(nightfury.Game)}, nil
			})

		games, err := nightfury.NewGamesFromRepo(repository)

		assert.NoError(t, err)
		if !cmp.Equal(expected, games) {
			assert.Fail(t, cmp.Diff(nightfury.Game{}, games))
		}
	})
}
