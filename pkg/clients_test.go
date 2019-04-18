package pkg_test

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/internal/mocks/db"
	"testing"
)

func TestClientAdd(t *testing.T) {
	t.Run("should be able to add game", func(t *testing.T) {
		client := pkg.NewClient("test", false)
		game := pkg.Game{
			Name:        "tic-tac-toe",
			Instruction: "instruction",
		}
		expected := pkg.Client{
			Name:         "test",
			GameStatuses: map[string]pkg.GameStatus{"tic-tac-toe": {Name: "tic-tac-toe", Status: "ready"}},
		}

		client.Add(game)

		if !cmp.Equal(expected, client) {
			assert.Fail(t, cmp.Diff(expected, client))
		}
	})

	t.Run("should replace the existing game", func(t *testing.T) {
		client := pkg.NewClient(
			"test",
			false,
			pkg.GameStatus{Name: "tic-tac-toe", Status: "started"},
		)
		game := pkg.Game{
			Name:        "tic-tac-toe",
			Instruction: "instruction",
		}
		expected := pkg.Client{
			Name:         "test",
			GameStatuses: map[string]pkg.GameStatus{"tic-tac-toe": {Name: "tic-tac-toe", Status: "ready"}},
		}

		client.Add(game)

		if !cmp.Equal(expected, client) {
			assert.Fail(t, cmp.Diff(expected, client))
		}
	})
}

func TestClientRemove(t *testing.T) {
	t.Run("should be able to remove a game", func(t *testing.T) {
		client := pkg.NewClient(
			"test",
			false,
			pkg.GameStatus{Name: "tic-tac-toe", Status: "started"},
		)
		expected := pkg.Client{
			Name:         "test",
			GameStatuses: map[string]pkg.GameStatus{},
		}

		client.Remove("tic-tac-toe")

		if !cmp.Equal(expected, client) {
			assert.Fail(t, cmp.Diff(expected, client))
		}
	})
}

func TestClientConnected(t *testing.T) {
	t.Run("should update the status", func(t *testing.T) {
		client := pkg.Client{Name: "client"}
		actual := client.Connected()
		expected := pkg.Client{Name: "client", Available: true}

		assert.Equal(t, expected, actual)
	})
}

func TestClientDisConnected(t *testing.T) {
	t.Run("should update the status", func(t *testing.T) {
		client := pkg.Client{Name: "client", Available: true}
		actual := client.Disconnected()
		expected := pkg.Client{Name: "client", Available: false}

		assert.Equal(t, expected, actual)
	})
}

func TestClientSave(t *testing.T) {
	t.Run("should be able to save client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := pkg.Client{Name: "client", Available: true}
		repository.EXPECT().Save("clients", client)

		err := client.Save(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := pkg.Client{Name: "client", Available: true}
		repository.EXPECT().Save("clients", client).Return(fmt.Errorf("unable to save"))

		err := client.Save(repository)

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
		repository.EXPECT().Fetch("clients", "one", gomock.Any()).Return(nil)

		actual, err := pkg.NewClientFromRepoWithName(repository, "one")

		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("should return error returned while fetching data from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("clients", "one", gomock.Any()).Return(fmt.Errorf("unable to fetch"))

		actual, err := pkg.NewClientFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch", err.Error())
		}
		assert.Equal(t, pkg.Client{}, actual)
	})
}

func TestNewClientsFromRepo(t *testing.T) {
	t.Run("should be able to get all the clients", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := pkg.Clients{"example": {Name: "example", Available: true, GameStatuses: pkg.GameStatuses{}}}
		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().FetchAll("clients", gomock.Any()).DoAndReturn(
			func(bucketName string, modelFn func(data []byte) (db.Model, error)) (interface{}, error) {
				data, _ := json.Marshal(pkg.Client{Name: "example", Available: true, GameStatuses: pkg.GameStatuses{}})
				model, err := modelFn(data)
				if err != nil {
					return nil, err
				}
				return pkg.Clients{model.ID(): model.(pkg.Client)}, nil
			})

		clients, err := pkg.NewClientsFromRepo(repository)

		assert.NoError(t, err)
		if !cmp.Equal(expected, clients) {
			assert.Fail(t, cmp.Diff(pkg.Client{}, clients))
		}
	})
}
