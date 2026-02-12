package domain_test

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"testing"
)

func TestProject_CreateSchema(t *testing.T) {
	project, err := domain.NewProject("test")
	require.NoError(t, err)

	t.Run("returns error if invalid name is provided", func(t *testing.T) {
		_, err := project.CreateSchema("", jsonschema.Schema{})
		require.ErrorIs(t, err, domain.ErrEmptySlug)
	})

	t.Run("returns error if invalid schema is provided", func(t *testing.T) {
		t.Skip()
		_, err := project.CreateSchema("user.created", jsonschema.Schema{})
		var perr errpack.Error
		require.ErrorAs(t, err, &perr)
		require.True(t, perr.Type() == errpack.TypeDomain)
	})
}
