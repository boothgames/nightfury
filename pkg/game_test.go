package pkg_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg"
	mocks "gitlab.com/jskswamy/nightfury/pkg/internal/mocks/db"
	"testing"
)

func TestGameSave(t *testing.T) {
	t.Run("should be able to save game", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		game := pkg.Game{Name: "game"}
		repository.EXPECT().Save("games", game)

		err := game.Save(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		game := pkg.Game{Name: "game"}
		repository.EXPECT().Save("games", game).Return(fmt.Errorf("unable to save"))

		err := game.Save(repository)

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
		repository.EXPECT().Fetch("games", "one", gomock.Any()).Return(nil)

		actual, err := pkg.NewGameFromRepoWithName(repository, "one")

		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("should return error returned while fetching data from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("games", "one", gomock.Any()).Return(fmt.Errorf("unable to fetch"))

		actual, err := pkg.NewGameFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch", err.Error())
		}
		assert.Equal(t, pkg.Game{}, actual)
	})
}
