package assert

import (
	"encoding/json"
	"fmt"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

// Game assert expected nightfury.Hint is the body of httptest.ResponseRecorder
func Game(t *testing.T, expected nightfury.Game, response *httptest.ResponseRecorder) {
	actual := nightfury.Game{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unable to unmarshal response as game, reason %v", err.Error()))
	}

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}

// Games assert expected nightfury.Hint is the body of httptest.ResponseRecorder
func Games(t *testing.T, expected nightfury.Games, response *httptest.ResponseRecorder) {
	actual := nightfury.Games{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unable to unmarshal response as games, reason %v", err.Error()))
	}

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}
