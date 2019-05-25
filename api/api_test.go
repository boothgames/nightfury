package api_test

import (
	"fmt"
	"github.com/boothgames/nightfury/api"
	"github.com/boothgames/nightfury/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testDBFileName = "test-db.db"

func TestPing(t *testing.T) {
	router := setupTestContext()
	defer teardownTestContext(t)
	response := performRequest(router, "GET", "/ping", nil)
	assert.Equal(t, http.StatusOK, response.Code)
}

func teardownTestContext(t *testing.T) {
	_ = db.Close()
	_ = os.RemoveAll(testDBFileName)
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	responseWriter := httptest.NewRecorder()
	r.ServeHTTP(responseWriter, req)
	return responseWriter
}

func setupTestContext() *gin.Engine {
	router := gin.Default()
	api.Bind(router)

	err := db.Initialize(testDBFileName)
	if err != nil {
		panic(fmt.Errorf("could not initialise test db %v", err))
	}
	return router
}
