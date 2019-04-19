package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var buckets = []string{"securityIncidents", "games"}

func TestPing(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)
	response := performRequest(router, "GET", "/ping", nil)
	assert.Equal(t, http.StatusOK, response.Code)
}

func teardownTestContext(t *testing.T) {
	for _, bucket := range buckets {
		if err := db.DeleteBucket(bucket); err != nil {
			panic(fmt.Errorf("error %v in deleting bucket %v", err, bucket))
		}
	}
	_ = db.Close()
}

func TestSecurityIncidentAPISuccessScenarios(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)
	t.Run("create security incident", func(t *testing.T) {
		expectedResponse := `{"Title":"title space title","Tag":"tag","Content":"new content","Takeaway":"new-takeaway2"}`

		createBody := `{"title": "title space title", "tag": "tag", "content": "new content", "takeaway": "new-takeaway2"}`
		response := performRequest(router, "POST", "/v1/security-incidents",
			bytes.NewBuffer([]byte(createBody)))

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
	t.Run("get all security incidents", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expectedResponse := fmt.Sprintf("{\"%s\":{\"Title\":\"%s\",\"Tag\":\"tag\",\"Content\":\"new content\",\"Takeaway\":\"new-takeaway2\"}}",
			titleHyphenated, title)

		response := performRequest(router, "GET", "/v1/security-incidents", nil)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
	t.Run("read security incident", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expectedResponse := fmt.Sprintf("{\"Title\":\"%s\",\"Tag\":\"tag\",\"Content\":\"new content\",\"Takeaway\":\"new-takeaway2\"}", title)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
	t.Run("update security incident", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expectedResponse := fmt.Sprintf("{\"Title\":\"%s\",\"Tag\":\"new tag\",\"Content\":\"new content\",\"Takeaway\":\"new-takeaway2\"}", title)

		updateBody := `{"title": "title space title", "tag": "new tag", "content": "new content", "takeaway": "new-takeaway2"}`
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated),
			bytes.NewBuffer([]byte(updateBody)))

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
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
		expectedResponse := fmt.Sprintf("{\"error\":\"securityIncident with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
}

func TestSecurityIncidentDeleteFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("delete security incident should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expectedResponse := fmt.Sprintf("{\"error\":\"securityIncident with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "DELETE", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
}

func TestSecurityIncidentUpdateFailure(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)

	t.Run("update security incident should fail when title does'nt exist in db", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expectedResponse := fmt.Sprintf("{\"error\":\"securityIncident with name %s doesn't exists\"}", titleHyphenated)

		response := performRequest(router, "GET", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated), nil)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
	t.Run("update security incident should fail if title is changed", func(t *testing.T) {
		title := "title space title"
		titleHyphenated := strings.Replace(title, " ", "-", -1)
		expectedResponse := fmt.Sprintf("{\"error\":\"title '%v' cannot be different\"}", title)

		createBody := `{"title": "title space title", "tag": "tag", "content": "new content", "takeaway": "new-takeaway2"}`
		performRequest(router, "POST", "/v1/security-incidents",
			bytes.NewBuffer([]byte(createBody)))

		updateBody := `{"title": "title space title2", "tag": "new tag", "content": "new content", "takeaway": "new-takeaway2"}`
		response := performRequest(router, "PUT", fmt.Sprintf("/v1/security-incidents/%v", titleHyphenated),
			bytes.NewBuffer([]byte(updateBody)))

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, expectedResponse, response.Body.String())
	})
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	responseWriter := httptest.NewRecorder()
	r.ServeHTTP(responseWriter, req)
	return responseWriter
}

func setupTestContext() *gin.Engine {
	router := gin.Default()
	Bind(router)
	err := db.Initialize("test-db.db")
	if err != nil {
		panic(fmt.Errorf("could not initialise test db %v", err))
	}
	return router
}
