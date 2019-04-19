package socket

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/olahol/melody.v1"
	"testing"
)

func Test_gameName(t *testing.T) {
	t.Run("it should return the game name from session", func(t *testing.T) {
		name, ok := gameName(&melody.Session{Keys: map[string]interface{}{"name": "game"}})

		assert.Equal(t, "game", name)
		assert.True(t, ok)
	})

	t.Run("it should return false game name is not available", func(t *testing.T) {
		name, ok := gameName(&melody.Session{Keys: map[string]interface{}{}})

		assert.Equal(t, "", name)
		assert.False(t, ok)
	})
}

func Test_clientID(t *testing.T) {
	t.Run("it should return the client id from session", func(t *testing.T) {
		id, ok := clientID(&melody.Session{Keys: map[string]interface{}{"id": "client"}})

		assert.Equal(t, "client", id)
		assert.True(t, ok)
	})

	t.Run("it should return false if client id is not available", func(t *testing.T) {
		id, ok := clientID(&melody.Session{Keys: map[string]interface{}{}})

		assert.Equal(t, "", id)
		assert.False(t, ok)
	})
}
