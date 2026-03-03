package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
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

	t.Run("try validation against schema", func(t *testing.T) {
		t.Skipf("just temporary impl")

		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		schema := jsonschema.Schema{}
		schema.Properties = make(map[string]*jsonschema.Schema)
		schema.Types = []string{"object"}

		name := &jsonschema.Schema{}
		name.Default = json.RawMessage(`"don"`)
		name.Types = []string{"string", "null"}

		age := &jsonschema.Schema{}
		age.Minimum = jsonschema.Ptr(1.0)
		age.Maximum = jsonschema.Ptr(100.0)
		age.Types = []string{"integer"}

		schema.Properties["name"] = name
		schema.Properties["age"] = age
		schema.Required = []string{"age", "name"}

		jData, err := schema.MarshalJSON()
		assert.NoError(t, err)
		t.Logf("%s", string(jData))

		resolved, err := schema.Resolve(nil)
		assert.NoError(t, err)
		user := User{
			Name: "",
			Age:  101,
		}

		uData, err := json.Marshal(user)
		assert.NoError(t, err)

		var input any
		err = json.Unmarshal(uData, &input)
		assert.NoError(t, err)

		err = resolved.Validate(input)
		assert.NoError(t, err)

		t.Logf("udata=%s", string(uData))
	})
}
