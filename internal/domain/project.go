package domain

import (
	"errors"

	"github.com/xplexer-lab/xplexer/internal/common/entity"
)

var (
	_ entity.Aggregate[ProjectState] = (*Project)(nil)
)

type Project struct {
	*entity.Entity
	schemas []string
}

type ProjectState struct {
	entity.State
	Schemas []string
}

func (p *Project) ToState() ProjectState {
	return ProjectState{}
}

func (p *Project) Load(s ProjectState) error {
	if err := p.Entity.Load(s.State); err != nil {
		return err
	}

	p.schemas = s.Schemas

	return nil
}

func NewProject() (*Project, error) {
	return nil, errors.New("not implemented")
}
