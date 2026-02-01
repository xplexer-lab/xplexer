package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectCreateIn struct {
	Name string `json:"name" validate:"required"`
}

type ProjectCreateOut struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var projectCreate = restapi.Operation(
	restapi.WithHandler(func(ctx context.Context, in ProjectCreateIn) (*ProjectCreateOut, error) {
		return consistency.Tx(ctx, func(ctx context.Context, repos Repos) (*ProjectCreateOut, error) {
			project, err := domain.NewProject(in.Name)

			if err != nil {
				return nil, err
			}

			if err := repos.Projects.Insert(ctx, project); err != nil {
				return nil, err
			}

			return nil, errpack.New("not implemented", errpack.Bootstrap())
		})
	}),
)
