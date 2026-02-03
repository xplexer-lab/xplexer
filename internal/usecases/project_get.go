package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectGetIn struct {
	ID string `path:"id"`
}

var projectGet = restapi.Operation(
	restapi.WithHandler(func(ctx context.Context, in ProjectGetIn) (*ProjectDto, error) {
		// todo: remove tx it is no needed here
		return consistency.Tx(ctx, func(ctx context.Context, repos Repositories) (*ProjectDto, error) {
			// todo: move id parsing into the binder
			id, err := entity.ParseId(in.ID)

			if err != nil {
				return nil, err
			}

			if project, err := repos.Projects.FindOne(ctx, id); err != nil {
				return nil, err
			} else {
				return newProjectDto(project), nil
			}
		})
	}),
)
