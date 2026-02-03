package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectCreateIn struct {
	Name string `json:"name" validate:"required"`
}

var projectCreate = restapi.Operation(
	restapi.WithHandler(func(ctx context.Context, in ProjectCreateIn) (*ProjectDto, error) {
		return consistency.Tx(ctx, func(ctx context.Context, repos Repositories) (*ProjectDto, error) {
			project, err := domain.NewProject(in.Name)

			if err != nil {
				return nil, err
			}

			if err := repos.Projects.Insert(ctx, project); err != nil {
				return nil, err
			}

			return newProjectDto(project), nil
		})
	}),
)
