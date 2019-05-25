package api_test

import (
	"fmt"
	internalAssert "github.com/boothgames/nightfury/api/internal/assert"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGameAPISuccessScenarios(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("create game incident", func(t *testing.T) {
		game := nightfury.Game{Name: "example", Instruction: "instruction", Type: "manual"}
		expected := nightfury.Game{Name: "example", Title: "", Instruction: "instruction", Type: "manual", Mode: "", Metadata: nil}

		response := performRequest(router, "POST", "/v1/games", game)

		assert.Equal(t, http.StatusCreated, response.Code)
		internalAssert.Game(t, expected, response)
	})

	t.Run("get all game incidents", func(t *testing.T) {
		expected := nightfury.Games{"example": {Name: "example", Title: "", Instruction: "instruction", Type: "manual", Mode: "", Metadata: nil}}

		response := performRequest(router, "GET", "/v1/games", nil)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.Games(t, expected, response)
	})

	t.Run("read game incident", func(t *testing.T) {
		gameName := "example"
		expected := nightfury.Game{Name: gameName, Title: "", Instruction: "instruction", Type: "manual", Mode: "", Metadata: nil}

		response := performRequest(router, "GET", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.Game(t, expected, response)
	})

	t.Run("update game incident", func(t *testing.T) {
		gameName := "example"

		expected := nightfury.Game{Name: "example", Title: "", Instruction: "new-instruction", Type: "manual", Mode: "", Metadata: nil}
		game := nightfury.Game{Name: "example", Title: "", Instruction: "new-instruction", Type: "manual", Mode: "", Metadata: nil}

		response := performRequest(router, "PUT", fmt.Sprintf("/v1/games/%v", gameName), game)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.Game(t, expected, response)
	})

	t.Run("delete game incident", func(t *testing.T) {
		gameName := "example"
		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestGameReadFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("read game should fail when name does'nt exist in db", func(t *testing.T) {
		gameName := "random"
		expected := fmt.Sprintf("{\"error\":\"game with name %v doesn't exists\"}", gameName)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestGameDeleteFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("delete hint should fail when name does'nt exist in db", func(t *testing.T) {
		gameName := "random"
		expected := fmt.Sprintf("{\"error\":\"game with name %v doesn't exists\"}", gameName)

		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestGameUpdateFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("update games should fail when name does'nt exist in db", func(t *testing.T) {
		gameName := "random"
		expected := fmt.Sprintf("{\"error\":\"game with name %v doesn't exists\"}", gameName)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})

	t.Run("update hint should fail if title is changed", func(t *testing.T) {

		name := "first"
		game := nightfury.Game{Name: name, Title: "", Instruction: "new-instruction", Type: "manual", Mode: "", Metadata: nil}

		performRequest(router, "POST", "/v1/games", game)

		expected := "{\"error\":\"name cannot be different\"}"

		updatedGame := nightfury.Game{Name: "updated-name", Title: "", Instruction: "new-instruction", Type: "manual", Mode: "", Metadata: nil}
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/games/%v", name), updatedGame)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}
