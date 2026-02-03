package usecases

import (
	"context"

	"github.com/samber/lo"
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/query"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectListIn struct{}

var projectList = restapi.Operation(
	restapi.WithOpMiddlewares(),
	restapi.WithHandler(func(ctx context.Context, in ProjectListIn) ([]ProjectDto, error) {
		repos := Repos(ctx)
		projects, err := repos.Projects.Find(ctx, query.And())

		if err != nil {
			return nil, err
		}

		return lo.Map(projects, func(p *domain.Project, _ int) ProjectDto {
			return *newProjectDto(p)
		}), nil

	}),
)
