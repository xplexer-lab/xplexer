package usecases

import (
	"github.com/xplexer-lab/xplexer/internal/common/entity"
	"github.com/xplexer-lab/xplexer/internal/domain"
)

type Repos struct {
	Project entity.Repository[*domain.Project, domain.ProjectState]
}
