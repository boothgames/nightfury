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

// SecurityIncident assert expected nightfury.SecurityIncident is the body of httptest.ResponseRecorder
func SecurityIncident(t *testing.T, expected nightfury.SecurityIncident, response *httptest.ResponseRecorder) {
	actual := nightfury.SecurityIncident{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unable to unmarshal response as security incident, reason %v", err.Error()))
	}

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}

// SecurityIncidents assert expected nightfury.SecurityIncident is the body of httptest.ResponseRecorder
func SecurityIncidents(t *testing.T, expected nightfury.SecurityIncidents, response *httptest.ResponseRecorder) {
	actual := nightfury.SecurityIncidents{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unable to unmarshal response as security incidents, reason %v", err.Error()))
	}

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}
