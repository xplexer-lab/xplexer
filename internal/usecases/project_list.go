package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectListIn struct{}

var projectList = restapi.Operation(
	restapi.WithOpMiddlewares(),
	restapi.WithHandler(func(ctx context.Context, in ProjectListIn) ([]ProjectDto, error) {
		return []ProjectDto{}, nil
	}),
)
