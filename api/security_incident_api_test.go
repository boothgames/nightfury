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

func TestSecurityIncidentAPISuccessScenarios(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("create security incident", func(t *testing.T) {
		incident := nightfury.SecurityIncident{Title: "title space title", Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		expected := nightfury.SecurityIncident{Title: "title space title", Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}

		response := performRequest(router, "POST", "/v1/security-incidents", incident)

		assert.Equal(t, http.StatusCreated, response.Code)
		internalAssert.SecurityIncident(t, expected, response)
	})

	t.Run("get all security incidents", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := nightfury.SecurityIncidents{titleHyphenated: {Title: title, Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}}

		response := performRequest(router, "GET", "/v1/security-incidents", nil)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.SecurityIncidents(t, expected, response)
	})

	t.Run("read security incident", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)

		expected := nightfury.SecurityIncident{Title: title, Tag: []string{"tag"}, Content: "new content", Takeaway: "new-takeaway2"}

		response := performRequest(router, "GET", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.SecurityIncident(t, expected, response)
	})

	t.Run("update security incident", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)

		incident := nightfury.SecurityIncident{Title: title, Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		expected := nightfury.SecurityIncident{Title: title, Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}

		response := performRequest(router, "PUT", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), incident)

		assert.Equal(t, http.StatusOK, response.Code)
		internalAssert.SecurityIncident(t, expected, response)
	})

	t.Run("delete security incident", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestSecurityIncidentReadFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("read security incident should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"securityIncident with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestSecurityIncidentDeleteFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("delete security incident should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"securityIncident with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestSecurityIncidentUpdateFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("update security incident should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"securityIncident with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})

	t.Run("update security incident should fail if title is changed", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expected := fmt.Sprintf("{\"error\":\"title '%v' cannot be different\"}", title)

		incident := nightfury.SecurityIncident{Title: title, Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		performRequest(router, "POST", "/v1/security-incidents", incident)

		updatedIncident := nightfury.SecurityIncident{Title: "title space title2", Tag: []string{"new tag"}, Content: "new content", Takeaway: "new-takeaway2"}
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), updatedIncident)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, expected, response.Body.String())
	})
}
