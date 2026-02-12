package domain

import (
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
	"strings"
)

var (
	_ entity.Aggregate[SchemaState] = new(Schema)
)

type SchemaId string

type SchemaSlug string

func (ss SchemaSlug) Validate() error {
	if len(strings.TrimSpace(string(ss))) == 0 {
		return ErrEmptySlug
	}

	return nil
}

type Schema struct {
	*entity.Entity
	slug    SchemaSlug
	project entity.Id
	version uint
}

func (sc *Schema) ToState() SchemaState {
	return SchemaState{
		State:   sc.Entity.ToState(),
		slug:    sc.slug,
		project: sc.project.Hex(),
	}
}

func (sc *Schema) Load(_ SchemaState) error {
	return ErrNotImplemented
}

type SchemaState struct {
	entity.State
	slug    SchemaSlug
	project string
}
