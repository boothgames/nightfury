package api_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGameAPISuccessScenarios(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)
	t.Run("create game incident", func(t *testing.T) {
		expectedResponse := `{"Name":"example","Instruction":"instruction","Type":"manual"}`

		createBody := `{"name": "example","instruction": "instruction","type": "manual"}`
		response := performRequest(router, "POST", "/v1/games",
			bytes.NewBuffer([]byte(createBody)))

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})

	t.Run("get all game incidents", func(t *testing.T) {
		expectedResponse := "{\"example\":{\"Name\":\"example\",\"Instruction\":\"instruction\",\"Type\":\"manual\"}}"

		response := performRequest(router, "GET", "/v1/games", nil)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})

	t.Run("read game incident", func(t *testing.T) {
		gameName := "example"
		expectedResponse := fmt.Sprintf("{\"Name\":\"%v\",\"Instruction\":\"instruction\",\"Type\":\"manual\"}", gameName)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})

	t.Run("update game incident", func(t *testing.T) {
		gameName := "example"
		expectedResponse := fmt.Sprintf("{\"Name\":\"%v\",\"Instruction\":\"new-instruction\",\"Type\":\"manual\"}", gameName)

		updateBody := "{\"Name\":\"example\",\"Instruction\":\"new-instruction\",\"Type\":\"manual\"}"
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/games/%v", gameName),
			bytes.NewBuffer([]byte(updateBody)))

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
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
		expectedResponse := fmt.Sprintf("{\"error\":\"game with name %v doesn't exists\"}", gameName)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
}

func TestGameDeleteFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("delete security incident should fail when name does'nt exist in db", func(t *testing.T) {
		gameName := "random"
		expectedResponse := fmt.Sprintf("{\"error\":\"game with name %v doesn't exists\"}", gameName)

		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
}

func TestGameUpdateFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("update games should fail when name does'nt exist in db", func(t *testing.T) {
		gameName := "random"
		expectedResponse := fmt.Sprintf("{\"error\":\"game with name %v doesn't exists\"}", gameName)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/games/%v", gameName), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
	t.Run("update security incident should fail if title is changed", func(t *testing.T) {
		firstName := "first"
		createBody := fmt.Sprintf("{\"name\": \"%v\",\"instruction\": \"instruction\",\"type\": \"manual\"}", firstName)
		performRequest(router, "POST", "/v1/games",
			bytes.NewBuffer([]byte(createBody)))

		expectedResponse := "{\"error\":\"name cannot be different\"}"

		updateBody := `{"name": "example","instruction": "instruction","type": "manual"}`
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/games/%v", firstName),
			bytes.NewBuffer([]byte(updateBody)))

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
}
