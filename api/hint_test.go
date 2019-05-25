package api_test

import (
	"fmt"
	internalAssert "github.com/boothgames/nightfury/api/internal/assert"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestHintAPISuccessScenarios(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("create hint", func(t *testing.T) {
		hint := nightfury.Hint{Title: "title space title", Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		expected := nightfury.Hint{Title: "title space title", Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}

		response := performRequest(router, "POST", "/v1/hints", hint)

		assert.Equal(t, http.StatusCreated, response.Code)
		internalAssert.Hint(t, expected, response)
	})

	t.Run("get all hints", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := nightfury.Hints{titleHyphenated: {Title: title, Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}}

		response := performRequest(router, "GET", "/v1/hints", nil)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.Hints(t, expected, response)
	})

	t.Run("read hint", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)

		expected := nightfury.Hint{Title: title, Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}

		response := performRequest(router, "GET", fmt.Sprintf("/v1/hints/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.Hint(t, expected, response)
	})

	t.Run("update hint", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)

		hint := nightfury.Hint{Title: title, Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		expected := nightfury.Hint{Title: title, Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}

		response := performRequest(router, "PUT", fmt.Sprintf("/v1/hints/%v", titleHyphenated), hint)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.Hint(t, expected, response)
	})

	t.Run("delete hint", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/hints/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestHintReadFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("read hint should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"hint with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/hints/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestHintDeleteFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("delete hint should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"hint with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/hints/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestHintUpdateFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("update hint should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"hint with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/hints/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})

	t.Run("update hint should fail if title is changed", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"title '%v' cannot be different\"}", title)

		hint := nightfury.Hint{Title: title, Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		performRequest(router, "POST", "/v1/hints", hint)

		updatedHint := nightfury.Hint{Title: "title space title2", Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/hints/%v", titleHyphenated), updatedHint)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}
