package api

import (
	"github.com/magiconair/properties/assert"
	"gopkg.in/olahol/melody.v1"
	"testing"
)

func Test_gameName(t *testing.T) {
	t.Run("it should return the game name from session", func(t *testing.T) {
		name := gameName(&melody.Session{Keys: map[string]interface{}{"name": "game"}})

		assert.Equal(t, "game", name)
	})
}
