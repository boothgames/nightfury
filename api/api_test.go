package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"io"
	"net/http"
	"net/http/httptest"
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

