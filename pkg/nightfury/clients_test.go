package nightfury_test

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/internal/mocks/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"testing"
)

func TestClientAdd(t *testing.T) {
	t.Run("should be able to add game", func(t *testing.T) {
		client := nightfury.NewClient("test", false)
		game := nightfury.Game{
			Name:        "tic-tac-toe",
			Instruction: "instruction",
		}
		expected := nightfury.Client{
			Name:         "test",
			GameStatuses: map[string]nightfury.GameStatus{"tic-tac-toe": {Name: "tic-tac-toe", Status: "ready"}},
		}

		client.Add(game)

		if !cmp.Equal(expected, client) {
			assert.Fail(t, cmp.Diff(expected, client))
		}
	})

	t.Run("should replace the existing game", func(t *testing.T) {
		client := nightfury.NewClient(
			"test",
			false,
			nightfury.GameStatus{Name: "tic-tac-toe", Status: "started"},
		)
		game := nightfury.Game{
			Name:        "tic-tac-toe",
			Instruction: "instruction",
		}
		expected := nightfury.Client{
			Name:         "test",
			GameStatuses: map[string]nightfury.GameStatus{"tic-tac-toe": {Name: "tic-tac-toe", Status: "ready"}},
		}

		client.Add(game)

		if !cmp.Equal(expected, client) {
			assert.Fail(t, cmp.Diff(expected, client))
		}
	})
}

func TestClientRemove(t *testing.T) {
	t.Run("should be able to remove a game", func(t *testing.T) {
		client := nightfury.NewClient(
			"test",
			false,
			nightfury.GameStatus{Name: "tic-tac-toe", Status: "started"},
		)
		expected := nightfury.Client{
			Name:         "test",
			GameStatuses: map[string]nightfury.GameStatus{},
		}

		client.Remove("tic-tac-toe")

		if !cmp.Equal(expected, client) {
			assert.Fail(t, cmp.Diff(expected, client))
		}
	})
}

func TestClientConnected(t *testing.T) {
	t.Run("should update the status", func(t *testing.T) {
		client := nightfury.Client{Name: "client"}
		actual := client.Connected()
		expected := nightfury.Client{Name: "client", Available: true}

		assert.Equal(t, expected, actual)
	})
}

func TestClientDisConnected(t *testing.T) {
	t.Run("should update the status", func(t *testing.T) {
		client := nightfury.Client{Name: "client", Available: true}
		actual := client.Disconnected()
		expected := nightfury.Client{Name: "client", Available: false}

		assert.Equal(t, expected, actual)
	})
}

func TestClientSave(t *testing.T) {
	t.Run("should be able to save client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := nightfury.Client{Name: "client", Available: true}
		repository.EXPECT().Save("clients", client)

		err := client.Save(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := nightfury.Client{Name: "client", Available: true}
		repository.EXPECT().Save("clients", client).Return(fmt.Errorf("unable to save"))

		err := client.Save(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestClientDelete(t *testing.T) {
	t.Run("should be able to save client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := nightfury.Client{Name: "client", Available: true}
		repository.EXPECT().Delete("clients", client)

		err := client.Delete(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := nightfury.Client{Name: "client", Available: true}
		repository.EXPECT().Delete("clients", client).Return(fmt.Errorf("unable to save"))

		err := client.Delete(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestNewClientFromRepoWithName(t *testing.T) {
	t.Run("should fetch the client from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("clients", "one", gomock.Any()).Return(true, nil)

		actual, err := nightfury.NewClientFromRepoWithName(repository, "one")

		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("should return error returned while fetching data from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("clients", "one", gomock.Any()).Return(false, fmt.Errorf("unable to fetch"))

		actual, err := nightfury.NewClientFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch", err.Error())
		}
		assert.Equal(t, nightfury.Client{}, actual)
	})
}

func TestNewClientsFromRepo(t *testing.T) {
	t.Run("should be able to get all the clients", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := nightfury.Clients{"example": {Name: "example", Available: true, GameStatuses: nightfury.GameStatuses{}}}
		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().FetchAll("clients", gomock.Any()).DoAndReturn(
			func(bucketName string, modelFn func(data []byte) (db.Model, error)) (interface{}, error) {
				data, _ := json.Marshal(nightfury.Client{Name: "example", Available: true, GameStatuses: nightfury.GameStatuses{}})
				model, err := modelFn(data)
				if err != nil {
					return nil, err
				}
				return nightfury.Clients{model.ID(): model.(nightfury.Client)}, nil
			})

		clients, err := nightfury.NewClientsFromRepo(repository)

		assert.NoError(t, err)
		if !cmp.Equal(expected, clients) {
			assert.Fail(t, cmp.Diff(nightfury.Client{}, clients))
		}
	})
}
