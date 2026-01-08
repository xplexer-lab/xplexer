package errpack_test

import (
	"errors"
	"fmt"
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

		t.Run("typed errors", func(t *testing.T) {
			testCases := []struct {
				name     string
				opt      errpack.Opt
				expected errpack.Type
			}{
				{"unknown", nil, errpack.Unknown},
				{"domain", errpack.WithDomain(), errpack.Domain},
				{"infra", errpack.WithInfra(), errpack.Infra},
			}

			for _, tt := range testCases {
				t.Run(tt.name, func(t *testing.T) {
					var err = errpack.New("error", tt.opt)
					assert.Equal(t, tt.expected, err.Type())
				})
			}
		})
	})

	t.Run("Unwrap", func(t *testing.T) {
		testCases := []struct {
			title string
			wrap  func(err error, msg string) error
		}{
			{"Wrap", func(err error, msg string) error {
				return errpack.Wrap(err, msg)
			}},
			{"WithPrev", func(err error, msg string) error {
				return errpack.New(msg, errpack.WithPrev(err))
			}},
		}

		for _, tt := range testCases {
			t.Run(tt.title, func(t *testing.T) {
				var err1 = errors.New("first")
				var err2 = fmt.Errorf("second: %w", err1)

				var wrapped = tt.wrap(err2, "wrapped")

				assert.True(t, errors.Is(wrapped, err1))
				assert.True(t, errors.Is(wrapped, err2))
			})
		}
	})

	t.Run("Wrap", func(t *testing.T) {
		t.Run("returns nil if empty item provided", func(t *testing.T) {
			var err = new(errpack.Error)

			assert.Nil(t, errpack.Wrap(nil, "wrapped"))
			assert.ErrorAs(
				t,
				errpack.Wrap(errors.New("base"), "wrapped"),
				&err,
			)
		})
	})
}
