package schema_test

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/xplexer-lab/xplexer/pkg/xkit/schema"
)

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func TestNewObject(t *testing.T) {
	o := schema.Object(
		schema.Property("first_name", schema.String(func(_ *jsonschema.Schema) {})),
		schema.Property("last_name", schema.String(func(_ *jsonschema.Schema) {})),
		schema.PropertyRequired("dick_size", schema.String(func(_ *jsonschema.Schema) {})),
		schema.PropertyRequired("ololo", schema.String(func(_ *jsonschema.Schema) {})),
		schema.Property2("some_property", func(b *schema.Property2Builder) {
			b.String().Required()
		}),
	)

	jData := must(o.MarshalJSON())

	t.Logf("%+v", o.String())
	t.Logf("%+v", string(jData))
}
