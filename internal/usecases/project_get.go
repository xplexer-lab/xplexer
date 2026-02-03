package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectGetIn struct {
	ID entity.Id `path:"id"`
}

var projectGet = restapi.Operation(
	restapi.WithHandler(func(ctx context.Context, in ProjectGetIn) (*ProjectDto, error) {
		repos := getRepos(ctx)

		if project, err := repos.Projects.FindOne(ctx, in.ID); err != nil {
			return nil, err
		} else {
			return newProjectDto(project), nil
		}
	}),
)
