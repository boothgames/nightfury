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

// Hint assert expected nightfury.Hint is the body of httptest.ResponseRecorder
func Hint(t *testing.T, expected nightfury.Hint, response *httptest.ResponseRecorder) {
	actual := nightfury.Hint{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unable to unmarshal response as hint, reason %v", err.Error()))
	}

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}

// Hints assert expected nightfury.Hint is the body of httptest.ResponseRecorder
func Hints(t *testing.T, expected nightfury.Hints, response *httptest.ResponseRecorder) {
	actual := nightfury.Hints{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unable to unmarshal response as hints, reason %v", err.Error()))
	}

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}
