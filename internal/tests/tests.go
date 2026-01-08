package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Tests struct {
	}
)

// Acceptance tests
func New() *Tests {
	return &Tests{}
}

func (tt *Tests) Run(t *testing.T) {
	t.Run("it works as expected", func(t *testing.T) {
		assert.True(t, false)
	})
}
