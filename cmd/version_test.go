package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildVersion(t *testing.T) {
	actual := BuildVersion()

	assert.Equal(t, "1.0-dev", actual)
}
