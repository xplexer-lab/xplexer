package domain

import (
	"fmt"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
)

var (
	_ entity.Aggregate[ProjectState] = (*Project)(nil)
)

func NewProject(name string) (*Project, error) {
	return &Project{
		Entity: entity.New(),
		name:   name,
	}, nil
}

func MustNewProject(name string) *Project {
	p, err := NewProject(name)

	if err != nil {
		panic(fmt.Errorf("must new project failed: %w", err))
	}

	return p
}

type Project struct {
	*entity.Entity
	name string
	// todo: extend with strategies
}

type ProjectState struct {
	entity.State
	Name string
}

type CreateSchemaCfg struct{}
type CreateSchemaOpt func(*CreateSchemaCfg)

func (p *Project) Name() string {
	return p.name
}

func (p *Project) ToState() ProjectState {
	return ProjectState{
		State: p.Entity.ToState(),
		Name:  p.name,
	}
}

func (p *Project) Load(s ProjectState) error {
	if err := p.Entity.Load(s.State); err != nil {
		return err
	}

	p.name = s.Name

	return nil
}

// CreateSchema create new schema
func (p *Project) CreateSchema(
	slug SchemaSlug,
	jSchema jsonschema.Schema,
	opts ...CreateSchemaOpt,
) (*Schema, error) {
	cfg := CreateSchemaCfg{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := slug.Validate(); err != nil {
		return nil, err
	}

	_ = jSchema

	var schema = Schema{
		Entity:  entity.New(),
		slug:    slug,
		project: p.Id(),
		version: 0,
	}

	if _, err := jSchema.MarshalJSON(); err != nil {
		return nil, err
	}

	// jSchema.Schema
	// todo: with rollback

	return &schema, nil
}
