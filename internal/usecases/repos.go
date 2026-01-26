package usecases

import (
	"github.com/xplexer-lab/xplexer/internal/common/entity"
	"github.com/xplexer-lab/xplexer/internal/domain"
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
