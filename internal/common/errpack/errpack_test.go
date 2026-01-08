package errpack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

func TestErrpack(t *testing.T) {
	t.Run("#New", func(t *testing.T) {

		t.Run("creates new error with unknown type by default", func(t *testing.T) {
			var err = errpack.New("error")
			assert.Equal(t, errpack.Unknown, err.Type())
		})

		t.Run("sets error type", func(t *testing.T) {
			var err = errpack.New("error", errpack.WithDomain())
			assert.Equal(t, errpack.Domain, err.Type())
		})

	})
}
