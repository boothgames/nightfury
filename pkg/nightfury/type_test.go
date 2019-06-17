package nightfury_test

import (
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlug(t *testing.T) {
	t.Run("should generate slug", func(t *testing.T) {
		slug := nightfury.Slug("WELCOME hELLO_World")

		assert.Equal(t, "welcome-hello-world", slug)
	})
}
