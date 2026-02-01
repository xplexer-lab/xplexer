package domain

import (
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
)

var (
	_ entity.Aggregate[SchemaState] = new(Schema)
)

type SchemaId string

type SchemaSlug string

type Schema struct {
	*entity.Entity
	slug    SchemaSlug
	project entity.Id
}

func (sc *Schema) ToState() SchemaState {
	return SchemaState{
		State:   sc.Entity.ToState(),
		slug:    sc.slug,
		project: sc.project.Hex(),
	}
}

func (sc *Schema) Load(_ SchemaState) error {
	return errpack.New("not implemented", errpack.Bootstrap())
}

type SchemaState struct {
	entity.State
	slug    SchemaSlug
	project string
}
