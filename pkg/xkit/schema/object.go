package schema

import "github.com/google/jsonschema-go/jsonschema"

type Opt func(*jsonschema.Schema)

func Object(
	opts ...Opt,
) *jsonschema.Schema {
	var s jsonschema.Schema
	s.Type = "object"

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

type ObjectBuilder struct {
}

func ensureProperties(s *jsonschema.Schema) {
	if s.Properties != nil {
		return
	}

	s.Properties = make(map[string]*jsonschema.Schema)
}

func Property(name string, schema *jsonschema.Schema) Opt {
	return func(s *jsonschema.Schema) {
		ensureProperties(s)
		s.Properties[name] = schema
	}
}

type Property2Builder struct {
	name     string
	required bool
}

func (pb *Property2Builder) String() *Property2Builder {
	return pb
}

func (pb *Property2Builder) SetRequired(required bool) *Property2Builder {
	pb.required = required
	return pb
}

func (pb *Property2Builder) Required() *Property2Builder {
	return pb.SetRequired(true)
}

func (pb *Property2Builder) Optional() *Property2Builder {
	return pb.SetRequired(false)
}

func (pb *Property2Builder) build() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
	}
}

func Property2(name string, f func(b *Property2Builder)) Opt {
	return func(s *jsonschema.Schema) {
		pb := &Property2Builder{name: name}
		f(pb)
		s.Properties[name] = pb.build()
		if pb.required {
			s.Required = append(s.Required, name)
		}
	}
}

func PropertyRequired(name string, schema *jsonschema.Schema) Opt {
	return func(s *jsonschema.Schema) {
		ensureProperties(s)
		s.Properties[name] = schema
		s.Required = append(s.Required, name)
	}
}

func String(f func(*jsonschema.Schema)) *jsonschema.Schema {
	var s jsonschema.Schema
	// todo: extract other
	s.Type = "string"
	f(&s)
	return &s
}

func OneOf(schemas ...*jsonschema.Schema) Opt {
	return func(s *jsonschema.Schema) {
		// todo: implement one of
		s.Types = []string{}
	}
}
