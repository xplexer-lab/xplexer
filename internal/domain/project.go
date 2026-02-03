package domain

import (
	"fmt"
	"maps"

	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
)

var (
	_ entity.Aggregate[ProjectState] = (*Project)(nil)
)

func NewProject(name string) (*Project, error) {
	return &Project{
		Entity:  entity.New(),
		name:    name,
		schemas: make(map[SchemaSlug]Schema),
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
	name    string
	schemas map[SchemaSlug]Schema
}

func (p *Project) Name() string {
	return p.name
}

type ProjectState struct {
	entity.State
	Name    string
	Schemas map[SchemaSlug]Schema
}

func (p *Project) ToState() ProjectState {
	return ProjectState{
		State:   p.Entity.ToState(),
		Name:    p.name,
		Schemas: maps.Clone(p.schemas),
	}
}

func (p *Project) Load(s ProjectState) error {
	if err := p.Entity.Load(s.State); err != nil {
		return err
	}

	p.name = s.Name
	p.schemas = s.Schemas

	return nil
}
