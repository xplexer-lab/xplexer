package usecases

import (
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
)

type Projects interface {
	entity.Repository[*domain.Project, domain.ProjectState]
}

type Repos struct {
	Projects Projects
	Users    interface {
		entity.Repository[*domain.Project, domain.ProjectState]
	}
}
