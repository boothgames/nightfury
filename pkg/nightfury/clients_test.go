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
			GameStatuses: map[string]nightfury.GameStatus{"tic-tac-toe": {Name: "tic-tac-toe", Status: nightfury.Ready}},
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
			nightfury.GameStatus{Name: "tic-tac-toe", Status: nightfury.InProgress},
		)
		game := nightfury.Game{
			Name:        "tic-tac-toe",
			Instruction: "instruction",
		}
		expected := nightfury.Client{
			Name:         "test",
			GameStatuses: map[string]nightfury.GameStatus{"tic-tac-toe": {Name: "tic-tac-toe", Status: nightfury.Ready}},
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
			nightfury.GameStatus{Name: "tic-tac-toe", Status: nightfury.InProgress},
		)
		expected := nightfury.Client{
			Name:         "test",
			GameStatuses: map[string]nightfury.GameStatus{},
		}

		client.Remove(nightfury.Game{Name: "tic-tac-toe"})

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

func TestClientStatus(t *testing.T) {
	type StatusScenario struct {
		name   string
		client nightfury.Client
		status nightfury.Status
	}

	scenarios := []StatusScenario{
		{
			name: "should return Ready as status",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Ready},
					"ludo":             {Status: nightfury.Ready},
					"snake-and-ladder": {Status: nightfury.Ready},
				},
			},
			status: nightfury.Ready,
		},
		{
			name: "should return Completed as status",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Completed},
					"ludo":             {Status: nightfury.Completed},
					"snake-and-ladder": {Status: nightfury.Completed},
				},
			},
			status: nightfury.Completed,
		},
		{
			name: "should return Failed as status when first game fails",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Failed},
					"ludo":             {Status: nightfury.Ready},
					"snake-and-ladder": {Status: nightfury.Ready},
				},
			},
			status: nightfury.Failed,
		},
		{
			name: "should return Failed as status when second game fails",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Completed},
					"ludo":             {Status: nightfury.Failed},
					"snake-and-ladder": {Status: nightfury.Ready},
				},
			},
			status: nightfury.Failed,
		},
		{
			name: "should return Failed as status when third game fails",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Completed},
					"ludo":             {Status: nightfury.Failed},
					"snake-and-ladder": {Status: nightfury.Ready},
				},
			},
			status: nightfury.Failed,
		},
		{
			name: "should return InProgress as status when other games are ready",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.InProgress},
					"ludo":             {Status: nightfury.Ready},
					"snake-and-ladder": {Status: nightfury.Ready},
				},
			},
			status: nightfury.InProgress,
		},
		{
			name: "should return InProgress as status when other games are ready/completed",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Completed},
					"ludo":             {Status: nightfury.InProgress},
					"snake-and-ladder": {Status: nightfury.Ready},
				},
			},
			status: nightfury.InProgress,
		},
		{
			name: "should return InProgress as status when other games are completed",
			client: nightfury.Client{
				GameStatuses: nightfury.GameStatuses{
					"tic-tac-toe":      {Status: nightfury.Completed},
					"ludo":             {Status: nightfury.Completed},
					"snake-and-ladder": {Status: nightfury.InProgress},
				},
			},
			status: nightfury.InProgress,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actualStatus := scenario.client.Status()

			assert.Equal(t, scenario.status, actualStatus)
		})
	}
}

func TestClientStart(t *testing.T) {
	t.Run("should start client", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Ready},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		mockRepository.EXPECT().Fetch("games", gomock.Any(), gomock.Any()).DoAndReturn(
			func(bucketName string, name string, model db.Model) (bool, error) {
				if name == "tic-tac-toe" || name == "ludo" || name == "snake-and-ladder" {
					return true, nil
				}
				return false, nil
			})

		_, err := client.Start()
		assert.NoError(t, err)
	})

	t.Run("should not start client", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.InProgress},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		_, err := client.Start()
		assert.Error(t, err)
		assert.Equal(t, "game already started", err.Error())
	})
}

func TestClientNext(t *testing.T) {
	t.Run("should return the next game", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Completed},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		mockRepository.EXPECT().Fetch("games", gomock.Any(), gomock.Any()).DoAndReturn(
			func(bucketName string, name string, model db.Model) (bool, error) {
				if name == "ludo" || name == "snake-and-ladder" {
					return true, nil
				}
				return false, nil
			})

		_, err := client.Next()

		assert.NoError(t, err)
	})

	t.Run("should not return the next game if game is not yet started", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Ready},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		_, err := client.Next()

		assert.Error(t, err)
		assert.Equal(t, "game not yet started", err.Error())
	})

	t.Run("should not return the next game if game is in progress", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.InProgress},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		_, err := client.Next()

		assert.Error(t, err)
		assert.Equal(t, "game already in progress", err.Error())
	})

	t.Run("should not return the next game if game is completed", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Completed},
				"ludo":             {Status: nightfury.Completed},
				"snake-and-ladder": {Status: nightfury.Completed},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		_, err := client.Next()

		assert.Error(t, err)
		assert.Equal(t, "game completed", err.Error())
	})

	t.Run("should not return the next game if game is failed", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Completed},
				"ludo":             {Status: nightfury.Failed},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepository(ctrl)
		restore := db.ReplaceDefaultRepositoryWith(mockRepository)

		defer func() {
			ctrl.Finish()
			restore()
		}()

		_, err := client.Next()

		assert.Error(t, err)
		assert.Equal(t, "game failed", err.Error())
	})
}

func TestClientHasNext(t *testing.T) {
	t.Run("should return the next game", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Completed},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		assert.True(t, client.HasNext())
	})

	t.Run("should not return the next game if game is not yet started", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Ready},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		assert.True(t, client.HasNext())
	})

	t.Run("should not return the next game if game is in progress", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.InProgress},
				"ludo":             {Status: nightfury.Ready},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		assert.True(t, client.HasNext())
	})

	t.Run("should not return the next game if game is completed", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Completed},
				"ludo":             {Status: nightfury.Completed},
				"snake-and-ladder": {Status: nightfury.Completed},
			},
		}

		assert.False(t, client.HasNext())
	})

	t.Run("should not return the next game if game is failed", func(t *testing.T) {
		client := nightfury.Client{
			GameStatuses: nightfury.GameStatuses{
				"tic-tac-toe":      {Status: nightfury.Completed},
				"ludo":             {Status: nightfury.Failed},
				"snake-and-ladder": {Status: nightfury.Ready},
			},
		}

		assert.False(t, client.HasNext())
	})
}
