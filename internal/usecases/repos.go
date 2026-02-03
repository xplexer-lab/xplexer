package usecases

import (
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
)

type TxManager = consistency.TxManager[Repositories]

type ProjectRepository interface {
	entity.Repository[*domain.Project, domain.ProjectState]
}

// Repositories set of repositories bounded by transaction
type Repositories struct {
	Projects ProjectRepository
}

var getRepos = consistency.Repos[Repositories]

// func repos(ctx context.Context) Repositories {
// 	return consistency.Repos[Repositories](ctx)
// }
